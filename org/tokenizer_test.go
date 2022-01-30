package org_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Ladicle/org2html/org"
)

func TestTokenize(t *testing.T) {
	var tests = []struct {
		desc       string
		input      string
		wantTokens []org.Token
	}{
		{
			desc: "empty",
		},
		{
			desc:  "tokenized successfully",
			input: "    CLOSED: [2022-01-30 Sun 10:03] ",
			wantTokens: []org.Token{
				org.NewToken(org.KindAgenda, 1, []string{"CLOSED", "2022-01-30", "Sun", "10:03", ""})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tokens, err := org.Tokenize(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: err=%v", err)
			}
			if !reflect.DeepEqual(tokens, tt.wantTokens) {
				t.Errorf("unexpected tokens:\ngot=%#v\nwant=%#v", tokens, tt.wantTokens)
			}
		})
	}
}
