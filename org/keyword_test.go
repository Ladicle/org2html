package org_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	. "github.com/Ladicle/org2html/org"
)

func TestLexKeyword(t *testing.T) {
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
			desc:      "keyword",
			line:      "#+key: value",
			wantFlag:  true,
			wantToken: NewToken(KindKeyword, 1, []string{"key", "value"}),
		},
		{
			desc:      "keyword with option",
			line:      "#+key[optional]: value",
			wantFlag:  true,
			wantToken: NewToken(KindKeyword, 1, []string{"key[optional]", "value"}),
		},
		{
			desc:      "keyword with space",
			line:      "   #+key: value",
			wantFlag:  true,
			wantToken: NewToken(KindKeyword, 1, []string{"key", "value"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			token, flag := LexKeyword(tt.line)
			if flag != tt.wantFlag {
				t.Errorf("unexpected flag: got=%v, want=%v", flag, tt.wantFlag)
			}
			if !reflect.DeepEqual(token, tt.wantToken) {
				t.Errorf("unexpected token:\ngot=%#v\nwant=%#v", token, tt.wantToken)
			}
		})
	}
}

func TestParseKeyword(t *testing.T) {
	var tests = []struct {
		desc      string
		token     Token
		wantNode  Node
		wantError error
	}{
		{
			desc:      "no keyword",
			token:     NewToken(KindKeyword, 1, []string{}),
			wantError: errors.New("keyword token[0] does not have 2 values: got=0"),
		},
		{
			desc:      "empty keyword key",
			token:     NewToken(KindKeyword, 1, []string{"", "value"}),
			wantError: errors.New("keyword key is empty"),
		},
		{
			desc:     "valid keyword",
			token:    NewToken(KindKeyword, 1, []string{"key", "value"}),
			wantNode: Keyword{Key: "KEY", Value: "value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			parser := DefaultParser([]Token{tt.token})
			consumed, node, err := ParseKeyword(&parser, 0)
			if err != nil {
				if tt.wantError == nil || err.Error() != tt.wantError.Error() {
					t.Fatalf("unexpected error: err=%v, want=%v", err, tt.wantError)
				}
				return
			} else if tt.wantError != nil {
				t.Fatalf("expect error but not occurred: want=%v", tt.wantError)
			}
			if want := 1; consumed != want {
				t.Errorf("unexpected consumed: got=%v, want=%v", consumed, want)
			}
			if !reflect.DeepEqual(node, tt.wantNode) {
				t.Errorf("unexpected node:\ngot=%#v\nwant=%#v", node, tt.wantNode)
			}
		})
	}
}

func TestKeywordWriter(t *testing.T) {
	var tests = []struct {
		desc    string
		kwd     Keyword
		wantOut string
	}{
		{
			desc:    "normal keyword",
			kwd:     Keyword{Key: "key", Value: "value"},
			wantOut: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var (
				out bytes.Buffer
			)
			if err := tt.kwd.Write(&out); err != nil {
				t.Fatalf("unexpected error: err=%v", err)
			}
			if got := out.String(); got != tt.wantOut {
				t.Errorf("unexpected output: got=%v, want=%v", got, tt.wantOut)
			}
		})
	}
}
