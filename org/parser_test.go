package org_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	. "github.com/Ladicle/org2html/org"
)

var _ Node = testNode{}

type testNode struct{}

func (n testNode) Write(w io.Writer) error { return nil }

var (
	testKind     = TokenKind("kind")
	testParserFn = func(p *Parser, i int) (consumed int, node Node, err error) {
		return 1, testNode{}, nil
	}
)

func TestParse(t *testing.T) {
	var tests = []struct {
		desc      string
		tokens    []Token
		wantNodes []Node
		wantError error
	}{
		{
			desc: "no token",
		},
		{
			desc:      "only one token",
			tokens:    []Token{NewToken(testKind, 1, []string{})},
			wantNodes: []Node{testNode{}},
		},
		{
			desc: "multiple tokens",
			tokens: []Token{
				NewToken(testKind, 1, []string{}),
				NewToken(testKind, 1, []string{}),
			},
			wantNodes: []Node{testNode{}, testNode{}},
		},
		{
			desc:      "unknown",
			tokens:    []Token{NewToken("unknown", 1, nil)},
			wantError: errors.New("unknown token: kind=unknown"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			parser := NewParser(tt.tokens, map[TokenKind]ParseFn{testKind: testParserFn})
			nodes, err := parser.Parse()
			if err != nil {
				if tt.wantError == nil || err.Error() != tt.wantError.Error() {
					t.Fatalf("unexpected error: err=%v, want=%v", err, tt.wantError)
				}
				return
			} else if tt.wantError != nil {
				t.Errorf("expect error but not occurred: want=%v", tt.wantError)
			}
			if !reflect.DeepEqual(nodes, tt.wantNodes) {
				t.Errorf("unexpected nodes:\ngot=%#v\nwant=%#v", nodes, tt.wantNodes)
			}
		})
	}
}
