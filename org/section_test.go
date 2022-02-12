package org_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	. "github.com/Ladicle/org2html/org"
)

func TestLexSection(t *testing.T) {
	var tests = []struct {
		desc      string
		line      string
		wantToken Token
		wantFlag  bool
	}{
		{
			desc:      "empty line",
			line:      "",
			wantFlag:  true,
			wantToken: NewToken(KindText, 1, []string{""}),
		},
		{
			desc:      "section",
			line:      "this is test section",
			wantFlag:  true,
			wantToken: NewToken(KindText, 1, []string{"this is test section"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			token, flag := LexText(tt.line)
			if flag != tt.wantFlag {
				t.Errorf("unexpected flag: got=%v, want=%v", flag, tt.wantFlag)
			}
			if !reflect.DeepEqual(token, tt.wantToken) {
				t.Errorf("unexpected token:\ngot=%#v\nwant=%#v", token, tt.wantToken)
			}
		})
	}
}

func TestParseSection(t *testing.T) {
	var tests = []struct {
		desc         string
		tokens       []Token
		wantConsumed int
		wantNode     Node
		wantError    error
	}{
		{
			desc:      "no section",
			tokens:    []Token{NewToken(KindText, 1, []string{})},
			wantError: errors.New("section token[0] does not have any values"),
		},
		{
			desc:         "one line",
			tokens:       []Token{NewToken(KindText, 1, []string{"this is section"})},
			wantNode:     Section{Paragraphs: []string{"this is section"}},
			wantConsumed: 1,
		},
		{
			desc: "multiple line",
			tokens: []Token{
				NewToken(KindText, 1, []string{"line1."}),
				NewToken(KindText, 1, []string{"line2."}),
			},
			wantNode:     Section{Paragraphs: []string{"line1. line2."}},
			wantConsumed: 2,
		},
		{
			desc: "multiple paragraphs",
			tokens: []Token{
				NewToken(KindText, 1, []string{"paragraph1."}),
				NewToken(KindText, 1, []string{""}),
				NewToken(KindText, 1, []string{"paragraph2."}),
				NewToken(KindText, 1, []string{"..."}),
			},
			wantNode:     Section{Paragraphs: []string{"paragraph1.", "paragraph2. ..."}},
			wantConsumed: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			parser := DefaultParser(tt.tokens)
			consumed, node, err := ParseSection(&parser, 0)
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

func TestSectionWriter(t *testing.T) {
	var tests = []struct {
		desc       string
		paragraphs []string
		wantOut    string
	}{
		{
			desc:       "one paragraph",
			paragraphs: []string{"this is section"},
			wantOut:    "<p>this is section</p>\n",
		},
		{
			desc:       "multiple paragraph",
			paragraphs: []string{"this is section1", "section2"},
			wantOut:    "<p>this is section1</p>\n<p>section2</p>\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var (
				out bytes.Buffer
				c   = Section{Paragraphs: tt.paragraphs}
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
