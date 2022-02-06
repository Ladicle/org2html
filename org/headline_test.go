package org_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	. "github.com/Ladicle/org2html/org"
)

func TestLexHeadline(t *testing.T) {
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
			desc:      "invalid headline",
			line:      " * invalid headline",
			wantToken: Token{},
		},
		{
			desc:      "Lv1 headline",
			line:      "* this is test headline",
			wantFlag:  true,
			wantToken: NewToken(KindHeadline, 1, []string{"*", "this is test headline"}),
		},
		{
			desc:      "Lv2 headline",
			line:      "** this is test headline",
			wantFlag:  true,
			wantToken: NewToken(KindHeadline, 1, []string{"**", "this is test headline"}),
		},
		{
			desc:      "Lv3 headline with meta",
			line:      "*** DONE [#A] this is test headline                                    :test_tag1:@tag2:",
			wantFlag:  true,
			wantToken: NewToken(KindHeadline, 1, []string{"***", "DONE [#A] this is test headline                                    :test_tag1:@tag2:"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			token, flag := LexHeadline(tt.line)
			if flag != tt.wantFlag {
				t.Errorf("unexpected flag: got=%v, want=%v", flag, tt.wantFlag)
			}
			if !reflect.DeepEqual(token, tt.wantToken) {
				t.Errorf("unexpected token:\ngot=%#v\nwant=%#v", token, tt.wantToken)
			}
		})
	}
}

func TestParseHeadline(t *testing.T) {
	var tests = []struct {
		desc      string
		token     Token
		wantNode  Node
		wantError error
	}{
		{
			desc:      "no headline",
			token:     NewToken(KindHeadline, 1, []string{}),
			wantError: errors.New("headline token[0] does not have enough values"),
		},
		{
			desc:     "Lv1 headline",
			token:    NewToken(KindHeadline, 1, []string{"*", "this is test headline"}),
			wantNode: Headline{Starts: 1, Title: "this is test headline"},
		},
		{
			desc:     "Lv2 headline with meta",
			token:    NewToken(KindHeadline, 1, []string{"**", "DONE [#A] this is test headline                                    :test_tag1:@tag2:"}),
			wantNode: Headline{Starts: 2, Title: "this is test headline", Keyword: "DONE", Priority: "A", Tags: []string{"test_tag1", "@tag2"}},
		},
		{
			desc:     "Lv3 headline with keyword",
			token:    NewToken(KindHeadline, 1, []string{"***", "WAIT this is test headline"}),
			wantNode: Headline{Starts: 3, Title: "this is test headline", Keyword: "WAIT"},
		},
		{
			desc:     "Lv4 headline with priority",
			token:    NewToken(KindHeadline, 1, []string{"****", "[#P1] this is test headline"}),
			wantNode: Headline{Starts: 4, Title: "this is test headline", Priority: "P1"},
		},
		{
			desc:     "Lv1 headline with tag",
			token:    NewToken(KindHeadline, 1, []string{"*", "this is test headline                                    :test_tag1:@tag2:"}),
			wantNode: Headline{Starts: 1, Title: "this is test headline", Tags: []string{"test_tag1", "@tag2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			parser := DefaultParser([]Token{tt.token})
			consumed, node, err := ParseHeadline(&parser, 0)
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

func TestHeadlineWriter(t *testing.T) {
	var tests = []struct {
		desc     string
		headline Headline
		wantOut  string
	}{
		{
			desc:     "Lv1 headline",
			headline: Headline{Starts: 1, Title: "this is test headline"},
			wantOut:  "<h1 class=\"org-headline\">\nthis is test headline\n</h1>\n",
		},
		{
			desc:     "Lv2 headline with keyword",
			headline: Headline{Starts: 2, Title: "this is test headline", Keyword: "TODO"},
			wantOut: `<h2 class="org-headline">
<span class="hl-kwd kwd-todo">TODO</span>
this is test headline
</h2>
`,
		},
		{
			desc:     "Lv3 headline with priority",
			headline: Headline{Starts: 3, Title: "this is test headline", Priority: "S"},
			wantOut: `<h3 class="org-headline">
<span class="hl-pri pri-s">S</span>
this is test headline
</h3>
`,
		},
		{
			desc: "Lv4 headline with all meta",
			headline: Headline{
				Starts:   3,
				Title:    "this is test headline",
				Keyword:  "WAIT",
				Priority: "P1",
				Tags:     []string{"@tag"},
			},
			wantOut: `<h3 class="org-headline">
<span class="hl-kwd kwd-wait">WAIT</span>
<span class="hl-pri pri-p1">P1</span>
this is test headline
</h3>
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var out bytes.Buffer
			if err := tt.headline.Write(&out); err != nil {
				t.Fatalf("unexpected error: err=%v", err)
			}
			if got := out.String(); got != tt.wantOut {
				t.Errorf("unexpected output: got=%v, want=%v", got, tt.wantOut)
			}
		})
	}
}
