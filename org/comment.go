package org

import (
	"fmt"
	"io"
	"regexp"
)

var _ Node = Comment{}

type Comment struct {
	Message string
}

func (c Comment) Write(w io.Writer) error {
	if c.Message != "" {
		fmt.Fprintf(w, "<!-- %s -->\n", c.Message)
	}
	return nil
}

var commentRegexp = regexp.MustCompile(`^\s*#\s*(.*)`)

func LexComment(line string) (Token, bool) {
	if m := commentRegexp.FindStringSubmatch(line); m != nil {
		return NewToken(KindComment, 1, []string{m[1]}), true
	}
	return Token{}, false
}

func ParseComment(p *Parser, i int) (int, Node, error) {
	if len(p.tokens[i].vals) < 1 {
		return 0, nil, fmt.Errorf("comment token[%v] does not have any values", i)
	}
	msg := p.tokens[i].vals[0]
	if msg == "" {
		return 1, nil, nil
	}
	return 1, Comment{Message: msg}, nil
}
