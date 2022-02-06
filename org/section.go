package org

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

var _ Node = Section{}

type Section struct {
	Paragraphs []string
}

func (s Section) Write(w io.Writer) error {
	for i := range s.Paragraphs {
		fmt.Fprintf(w, "<p>%s</p>\n", s.Paragraphs[i])
	}
	return nil
}

func LexSection(line string) (Token, bool) {
	return NewToken(KindSection, 1, []string{strings.TrimSpace(line)}), true
}

func ParseSection(p *Parser, i int) (int, Node, error) {
	var (
		buf        bytes.Buffer
		paragraphs []string
	)
	start, end := i, len(p.tokens)
	for i < end && p.tokens[i].kind == KindSection {
		if len(p.tokens[i].vals) == 0 {
			return 0, nil, fmt.Errorf("section token[%d] does not have any values", i)
		}
		line := p.tokens[i].vals[0]
		i++
		// start new paragraph
		if line == "" {
			if para := buf.String(); para != "" {
				paragraphs = append(paragraphs, para)
			}
			buf.Reset()
			continue
		}
		if buf.Len() > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(line)
	}
	// pack rest paragraph
	if para := buf.String(); para != "" {
		paragraphs = append(paragraphs, para)
	}
	return i - start, Section{Paragraphs: paragraphs}, nil
}
