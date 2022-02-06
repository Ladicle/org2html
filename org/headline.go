package org

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

var _ Node = Headline{}

type Headline struct {
	Starts   int
	Title    string
	Keyword  string   // optional
	Priority string   // optional
	Tags     []string // optional

	// TODO: additional fields

	// properties
	// agenda
	// logbook

	// section
}

// Write writes headline data as HTML elements to the specified writer.
// NOTE: headline tags are ignored.
func (h Headline) Write(w io.Writer) error {
	if h.Starts == 0 {
		return fmt.Errorf("invalid number of starts: %#v", h)
	}
	if h.Title == "" {
		return fmt.Errorf("title is empty: %#v", h)
	}

	fmt.Fprintf(w, "<h%d class=\"org-headline\">\n", h.Starts) // TODO: add ID
	if h.Keyword != "" {
		fmt.Fprintf(w, "<span class=\"hl-kwd kwd-%s\">%s</span>\n", strings.ToLower(h.Keyword), h.Keyword)
	}
	if h.Priority != "" {
		fmt.Fprintf(w, "<span class=\"hl-pri pri-%s\">%s</span>\n", strings.ToLower(h.Priority), h.Priority)
	}
	fmt.Fprintln(w, h.Title)
	fmt.Fprintf(w, "</h%d>\n", h.Starts)
	return nil
}

var (
	headlineRegexp = regexp.MustCompile(`^([*]+)\s+(.*)`)
	hlDataRegexp   = regexp.MustCompile(`^(?:([A-Z]+)\s+)?(?:\[#(\w+)\]\s+)?(.*?)(?:\s+:([A-Za-z0-9_@#%:]+):\s*)?$`)
)

func LexHeadline(line string) (Token, bool) {
	if m := headlineRegexp.FindStringSubmatch(line); m != nil {
		return NewToken(KindHeadline, 1, m[1:]), true
	}
	return Token{}, false
}

func ParseHeadline(p *Parser, i int) (int, Node, error) {
	if len(p.tokens[i].vals) < 2 {
		return 0, nil, fmt.Errorf("headline token[%d] does not have enough values", i)
	}

	m := hlDataRegexp.FindStringSubmatch(p.tokens[i].vals[1])
	if m == nil {
		return 0, nil, fmt.Errorf("headline token[%d] has invalid format title: %s", i, p.tokens[i].vals[0])
	}

	hl := Headline{
		Starts:   len(p.tokens[i].vals[0]),
		Keyword:  m[1], // TODO: validation
		Priority: m[2],
		Title:    m[3],
	}
	if m[4] != "" {
		hl.Tags = strings.Split(m[4], ":")
	}

	return 1, hl, nil
}
