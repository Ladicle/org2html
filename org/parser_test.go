package org_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Ladicle/org2html/org"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		desc      string
		tokens    []org.Token
		wantNodes []org.Node
		wantError error
	}{
		{
			desc: "no token",
		},
		{
			desc: "only one token",
			tokens: []org.Token{
				org.NewToken(org.KindAgenda, 1, []string{"CLOSED", "2022-01-30", "Sun", "10:03", ""})},
			wantNodes: []org.Node{
				org.Agenda{Logs: map[org.AgendaKey]org.Timestamp{
					org.AgendaClosed: mustParseTimestamp(t, "2022-01-30 Sun 10:03", "")}}},
		},
		{
			desc: "multiple tokens",
			tokens: []org.Token{
				org.NewToken(org.KindAgenda, 1, []string{"CLOSED", "2022-01-30", "Sun", "10:03", ""}),
				org.NewToken(org.KindAgenda, 1, []string{"CLOSED", "2022-01-30", "Sun", "10:03", ""})},
			wantNodes: []org.Node{
				org.Agenda{Logs: map[org.AgendaKey]org.Timestamp{
					org.AgendaClosed: mustParseTimestamp(t, "2022-01-30 Sun 10:03", "")}},
				org.Agenda{Logs: map[org.AgendaKey]org.Timestamp{
					org.AgendaClosed: mustParseTimestamp(t, "2022-01-30 Sun 10:03", "")}}},
		},
		{
			desc:      "unknown",
			tokens:    []org.Token{org.NewToken("unknown", 1, nil)},
			wantError: errors.New("unknown token: kind=unknown"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			parser := org.NewParser(tt.tokens)
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
