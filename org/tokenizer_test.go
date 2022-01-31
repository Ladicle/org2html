package org_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	. "github.com/Ladicle/org2html/org"
)

var (
	testToken = NewToken("test", 0, []string{})
	testLexFn = func(line string) (t Token, ok bool) {
		if line == "test" {
			return testToken, true
		}
		return Token{}, false
	}
)

func TestTokenize(t *testing.T) {
	var tests = []struct {
		desc       string
		input      string
		wantTokens []Token
		wantError  error
	}{
		{
			desc: "empty",
		},
		{
			desc:       "one token",
			input:      "test",
			wantTokens: []Token{testToken},
		},
		{
			desc:       "multiple token",
			input:      "test\ntest",
			wantTokens: []Token{testToken, testToken},
		},
		{
			desc:      "no lexers can parse",
			input:     "invalid",
			wantError: errors.New("no lexers can parse \"invalid\""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tokenizer := NewTokenizer([]LexFn{testLexFn})
			tokens, err := tokenizer.Tokenize(strings.NewReader(tt.input))
			if err != nil {
				if tt.wantError == nil || err.Error() != tt.wantError.Error() {
					t.Fatalf("unexpected error: err=%v, want=%v", err, tt.wantError)
				}
				return
			} else if tt.wantError != nil {
				t.Fatalf("expect error but not occurred: want=%v", tt.wantError)
			}
			if !reflect.DeepEqual(tokens, tt.wantTokens) {
				t.Errorf("unexpected tokens:\ngot=%#v\nwant=%#v", tokens, tt.wantTokens)
			}
		})
	}
}
