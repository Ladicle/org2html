package org

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// TODO: details/summary

var _ Node = Block{}

type Block struct {
	Name    string
	Content string
}

var _ Node = SourceBlock{}

type SourceBlock struct {
	Language   string
	SourceCode string
	Property   []string
}

func (c Block) Write(w io.Writer) error {
	fmt.Fprintf(w, "<div class=\"org-block block-%s\">\n", strings.ToLower(c.Name))
	fmt.Fprintln(w, c.Content)
	fmt.Fprintln(w, "</div>")
	return nil
}

func (c SourceBlock) Write(w io.Writer) error {
	fmt.Fprintln(w, "<div class=\"org-block block-src\">")
	fmt.Fprintf(w, "<code class=\"block lang-%s\" data-lang=\"%s\">\n", c.Language, c.Language)
	fmt.Fprintln(w, c.SourceCode)
	fmt.Fprintln(w, "</code>")
	fmt.Fprintln(w, "</div>")
	return nil
}

const sourceBlockName = "SRC"

var (
	beginBlockRegexp = regexp.MustCompile(`(?i)^\s*#\+BEGIN(?:_(\w+))(?:\s+(.*))?`)
	endBlockRegexp   = regexp.MustCompile(`(?i)^\s*#\+END(?:_(\w+))`)
)

func LexBlock(line string) (Token, bool) {
	if m := beginBlockRegexp.FindStringSubmatch(line); m != nil {
		return NewToken(KindBlockBegin, 1, m[1:]), true
	} else if m := endBlockRegexp.FindStringSubmatch(line); m != nil {
		return NewToken(KindBlockEnd, 1, m[1:]), true
	}
	return Token{}, false
}

func ParseBlock(p *Parser, i int) (int, Node, error) {
	if got, want := len(p.tokens[i].vals), 2; got != want {
		return 0, nil, fmt.Errorf("block token[%d] does not have %d values: got=%d", i, want, got)
	}

	var block Block
	switch name := strings.ToUpper(p.tokens[i].vals[0]); name {
	case "":
		return 0, nil, errors.New("block name is empty")
	default:
		block = Block{Name: name}
	}

	var (
		content bytes.Buffer
		start   = i
	)
	for i++; i < len(p.tokens); i++ {
		if p.tokens[i].kind == KindBlockEnd {
			if got, want := strings.ToUpper(p.tokens[i].vals[0]), block.Name; got != want {
				return 0, nil, fmt.Errorf(
					"token[%d] is unexpected block end: got=%v, want=%v", i, got, want)
			}
			break
		}
		content.WriteString(p.tokens[i].vals[0])
		content.WriteString("\n")
	}
	block.Content = strings.TrimRight(content.String(), "\n")

	if block.Name != sourceBlockName {
		return i - start + 1, block, nil
	}

	// start extra parsing for source block.
	parts := strings.SplitN(p.tokens[start].vals[1], " ", 2)

	lang := parts[0]
	var property string
	if len(parts) == 2 {
		property = strings.TrimPrefix(strings.TrimSpace(parts[1]), ":")
	}

	var srcBlock = SourceBlock{
		Language:   lang,
		SourceCode: block.Content,
	}
	if property != "" {
		for _, v := range strings.Split(property, " :") {
			srcBlock.Property = append(srcBlock.Property, strings.TrimSpace(v))
		}
	}
	return i - start + 1, srcBlock, nil
}
