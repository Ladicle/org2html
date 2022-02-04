package org

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var _ Node = Keyword{}

type Keyword struct {
	Key   string
	Value string
}

func (k Keyword) Write(w io.Writer) error {
	// noop
	return nil
}

// org syntax supports optional, but I couldn't find out the use-cases.
// Therefore, the optional value ignore here.
// https://orgmode.org/worg/dev/org-syntax.html#Affiliated_keywords
var keywordRegexp = regexp.MustCompile(`^\s*#\+([^:]+):\s*(.*)`)

func LexKeyword(line string) (Token, bool) {
	if m := keywordRegexp.FindStringSubmatch(line); m != nil {
		return NewToken(KindKeyword, 1, m[1:]), true
	}
	return Token{}, false
}

func ParseKeyword(p *Parser, i int) (int, Node, error) {
	if len(p.tokens[i].vals) != 2 {
		return 0, nil, fmt.Errorf(
			"keyword token[%v] does not have 2 values: got=%v", i, len(p.tokens[i].vals))
	}

	var (
		key = strings.ToUpper(p.tokens[i].vals[0])
		val = strings.TrimSpace(p.tokens[i].vals[1])
	)
	switch KeywordType(key) {
	case "":
		return 0, nil, errors.New("keyword key is empty")
	default:
		return 1, Keyword{Key: key, Value: val}, nil
	}
}
