package org_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	. "github.com/Ladicle/org2html/org"
)

func TestLexBlock(t *testing.T) {
	var tests = []struct {
		desc      string
		line      string
		wantToken Token
		wantFlag  bool
	}{
		{
			desc:      "empty line",
			line:      "",
			wantToken: Token{},
		},
		{
			desc:      "block begin",
			line:      "#+begin_info",
			wantFlag:  true,
			wantToken: NewToken(KindBlockBegin, 1, []string{"info", ""}),
		},
		{
			desc:      "block begin (upper)",
			line:      "#+BEGIN_EXAMPLE",
			wantFlag:  true,
			wantToken: NewToken(KindBlockBegin, 1, []string{"EXAMPLE", ""}),
		},
		{
			desc:      "block begin with space",
			line:      "    #+begin_quote  ",
			wantFlag:  true,
			wantToken: NewToken(KindBlockBegin, 1, []string{"quote", ""}),
		},
		{
			desc:      "block begin with property",
			line:      "#+begin_src bash :details t",
			wantFlag:  true,
			wantToken: NewToken(KindBlockBegin, 1, []string{"src", "bash :details t"}),
		},
		{
			desc:      "block end",
			line:      "#+end_src",
			wantFlag:  true,
			wantToken: NewToken(KindBlockEnd, 1, []string{"src"}),
		},
		{
			desc:      "block end (upper)",
			line:      "#+END_QUOTE",
			wantFlag:  true,
			wantToken: NewToken(KindBlockEnd, 1, []string{"QUOTE"}),
		},
		{
			desc:      "block end with space",
			line:      "   #+end_src",
			wantFlag:  true,
			wantToken: NewToken(KindBlockEnd, 1, []string{"src"}),
		},
		{
			desc: "escaped block begin",
			line: ",#+begin_info",
		},
		{
			desc: "escaped block end",
			line: ",#+end_info",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			token, flag := LexBlock(tt.line)
			if flag != tt.wantFlag {
				t.Errorf("unexpected flag: got=%v, want=%v", flag, tt.wantFlag)
			}
			if !reflect.DeepEqual(token, tt.wantToken) {
				t.Errorf("unexpected token:\ngot=%#v\nwant=%#v", token, tt.wantToken)
			}
		})
	}
}

func TestParseBlock(t *testing.T) {
	var tests = []struct {
		desc         string
		tokens       []Token
		wantConsumed int
		wantNode     Node
		wantError    error
	}{
		{
			desc:      "empty block",
			tokens:    []Token{NewToken(KindBlockBegin, 1, []string{})},
			wantError: errors.New("block token[0] does not have 2 values: got=0"),
		},
		{
			desc:      "no name",
			tokens:    []Token{NewToken(KindBlockBegin, 1, []string{"", ""})},
			wantError: errors.New("block name is empty"),
		},
		{
			desc: "mismatch block",
			tokens: []Token{
				NewToken(KindBlockBegin, 1, []string{"info", ""}),
				NewToken(KindBlockEnd, 1, []string{"quote", ""}),
			},
			wantError: errors.New("token[1] is unexpected block end: got=QUOTE, want=INFO"),
		},
		{
			desc: "empty block",
			tokens: []Token{
				NewToken(KindBlockBegin, 1, []string{"info", ""}),
				NewToken(KindBlockEnd, 1, []string{"info", ""}),
			},
			wantNode:     Block{Name: "INFO", Content: ""},
			wantConsumed: 2,
		},
		{
			desc: "block with content",
			tokens: []Token{
				NewToken(KindBlockBegin, 1, []string{"QUOTE", ""}),
				NewToken(KindText, 1, []string{"hello"}),
				NewToken(KindText, 1, []string{"world"}),
				NewToken(KindBlockEnd, 1, []string{"QUOTE", ""}),
			},
			wantNode:     Block{Name: "QUOTE", Content: "hello\nworld"},
			wantConsumed: 4,
		},
		{
			desc: "source code block",
			tokens: []Token{
				NewToken(KindBlockBegin, 1, []string{"SRC", "bash"}),
				NewToken(KindText, 1, []string{"#!/bin/bash"}),
				NewToken(KindText, 1, []string{"set -ex"}),
				NewToken(KindText, 1, []string{"echo 'hello'"}),
				NewToken(KindBlockEnd, 1, []string{"SRC", ""}),
			},
			wantNode: SourceBlock{
				Language:   "bash",
				SourceCode: "#!/bin/bash\nset -ex\necho 'hello'",
			},
			wantConsumed: 5,
		},
		{
			desc: "source code block with property",
			tokens: []Token{
				NewToken(KindBlockBegin, 1, []string{"SRC", "go  :var x=2  :var y=3"}),
				NewToken(KindText, 1, []string{"package main"}),
				NewToken(KindText, 1, []string{""}),
				NewToken(KindText, 1, []string{"func main() { // noop }"}),
				NewToken(KindBlockEnd, 1, []string{"SRC", ""}),
			},
			wantNode: SourceBlock{
				Language:   "go",
				SourceCode: "package main\n\nfunc main() { // noop }",
				Property:   []string{"var x=2", "var y=3"},
			},
			wantConsumed: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			parser := DefaultParser(tt.tokens)
			consumed, node, err := ParseBlock(&parser, 0)
			if err != nil {
				if tt.wantError == nil || err.Error() != tt.wantError.Error() {
					t.Fatalf("unexpected error: err=%v, want=%v", err, tt.wantError)
				}
				return
			} else if tt.wantError != nil {
				t.Fatalf("expect error but not occurred: want=%v", tt.wantError)
			}
			if consumed != tt.wantConsumed {
				t.Errorf("unexpected consumed: got=%v, want=%v", consumed, tt.wantConsumed)
			}
			if !reflect.DeepEqual(node, tt.wantNode) {
				t.Errorf("unexpected node:\ngot=%#v\nwant=%#v", node, tt.wantNode)
			}
		})
	}
}

func TestBlockWriter(t *testing.T) {
	var tests = []struct {
		desc    string
		msg     string
		wantOut string
	}{
		{
			desc:    "normal block",
			msg:     "hello world",
			wantOut: "<!-- hello world -->\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var (
				out bytes.Buffer
				c   = Block{}
			)
			if err := c.Write(&out); err != nil {
				t.Fatalf("unexpected error: err=%v", err)
			}
			if got := out.String(); got != tt.wantOut {
				t.Errorf("unexpected output: got=%v, want=%v", got, tt.wantOut)
			}
		})
	}
}
