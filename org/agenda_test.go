package org_test

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/Ladicle/org2html/org"
)

func TestLexAgenda(t *testing.T) {
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
			desc:     "closed with timestamp",
			line:     "    CLOSED: [2022-01-30 Sun 10:03] ",
			wantFlag: true,
			wantToken: NewToken(KindAgenda, 1, []string{
				"CLOSED", "2022-01-30", "Sun", "10:03", ""}),
		},
		{
			desc:     "schedule with interval",
			line:     "  SCHEDULED: <2022-01-30 Sun +1w>",
			wantFlag: true,
			wantToken: NewToken(KindAgenda, 1, []string{
				"SCHEDULED", "2022-01-30", "Sun", "", "+1w"}),
		},
		{
			desc:     "deadline",
			line:     "DEADLINE: <2022-01-30 Sun>      ",
			wantFlag: true,
			wantToken: NewToken(KindAgenda, 1, []string{
				"DEADLINE", "2022-01-30", "Sun", "", ""}),
		},
		{
			desc:     "multiple agenda",
			line:     "DEADLINE: <2022-01-30 Sun>   SCHEDULED: <2022-01-30 Sun>",
			wantFlag: true,
			wantToken: NewToken(KindAgenda, 2, []string{
				"DEADLINE", "2022-01-30", "Sun", "", "",
				"SCHEDULED", "2022-01-30", "Sun", "", ""}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			token, flag := LexAgenda(tt.line)
			if flag != tt.wantFlag {
				t.Errorf("unexpected flag: got=%v, want=%v", flag, tt.wantFlag)
			}
			if !reflect.DeepEqual(token, tt.wantToken) {
				t.Errorf("unexpected token:\ngot=%#v\nwant=%#v", token, tt.wantToken)
			}
		})
	}
}

func TestParseAgenda(t *testing.T) {
	var tests = []struct {
		desc      string
		token     Token
		wantNode  Node
		wantError error
	}{
		{
			desc: "one item",
			token: NewToken(KindAgenda, 1, []string{
				"CLOSED", "2022-01-30", "Sun", "10:03", ""}),
			wantNode: Agenda{Logs: map[AgendaKey]Timestamp{
				AgendaClosed: mustParseTimestamp(t, "2022-01-30 Sun 10:03", "")},
			},
		},
		{
			desc: "multiple items",
			token: NewToken(KindAgenda, 2, []string{
				"DEADLINE", "2022-01-30", "Sun", "", "",
				"SCHEDULED", "2022-01-30", "Sun", "", "+1w"}),
			wantNode: Agenda{Logs: map[AgendaKey]Timestamp{
				AgendaDeadline:  mustParseDatestamp(t, "2022-01-30 Sun", ""),
				AgendaScheduled: mustParseDatestamp(t, "2022-01-30 Sun", "+1w")}},
		},
		{
			desc:      "out of range",
			token:     NewToken(KindAgenda, 1, []string{}),
			wantError: errors.New("agenda item number and its values are unmatched: num=1, vals=[]string{}"),
		},
		{
			desc: "invalid",
			token: NewToken(KindAgenda, 1, []string{
				"CLOSED", "2022-01-30", "Invalid", "10:04", "++2d"}),
			wantNode: Agenda{Logs: map[AgendaKey]Timestamp{
				AgendaClosed: mustParseTimestamp(t, "2022-01-30 Sun 10:04", "++2d")},
			},
			wantError: errors.New("parsing time \"2022-01-30 Invalid 10:04\" as \"2006-01-02 Mon 15:04\": cannot parse \"Invalid 10:04\" as \"Mon\""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			parser := DefaultParser([]Token{tt.token})
			consumed, node, err := ParseAgenda(&parser, 0)
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
