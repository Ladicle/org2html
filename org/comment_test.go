package org_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	. "github.com/Ladicle/org2html/org"
)

func TestLexComment(t *testing.T) {
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
			desc:      "comment",
			line:      "# this is test comment",
			wantFlag:  true,
			wantToken: NewToken(KindComment, 1, []string{"this is test comment"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			token, flag := LexComment(tt.line)
			if flag != tt.wantFlag {
				t.Errorf("unexpected flag: got=%v, want=%v", flag, tt.wantFlag)
			}
			if !reflect.DeepEqual(token, tt.wantToken) {
				t.Errorf("unexpected token:\ngot=%#v\nwant=%#v", token, tt.wantToken)
			}
		})
	}
}

func TestParseComment(t *testing.T) {
	var tests = []struct {
		desc      string
		token     Token
		wantNode  Node
		wantError error
	}{
		{
			desc:      "no comment",
			token:     NewToken(KindComment, 1, []string{}),
			wantError: errors.New("comment token[0] does not have any values"),
		},
		{
			desc:  "empty comment",
			token: NewToken(KindComment, 1, []string{""}),
		},
		{
			desc:     "valid comment",
			token:    NewToken(KindComment, 1, []string{"this is comment"}),
			wantNode: Comment{Message: "this is comment"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			parser := DefaultParser([]Token{tt.token})
			consumed, node, err := ParseComment(&parser, 0)
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

func TestCommentWriter(t *testing.T) {
	var tests = []struct {
		desc    string
		msg     string
		wantOut string
	}{
		{
			desc:    "normal comment",
			msg:     "hello world",
			wantOut: "<!-- hello world -->",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var (
				out bytes.Buffer
				c   = Comment{Message: tt.msg}
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
