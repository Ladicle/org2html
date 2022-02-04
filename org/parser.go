package org

import (
	"fmt"
)

// ParseFn is a function to parse Node and return the number of tokens consumed,
// parsed Node and error. The argument i indicates the index of the token that
// will start parsing.
type ParseFn = func(p *Parser, i int) (consumed int, node Node, err error)

// defaultParseFns expresses currently supported passers.
var defaultParseFns = map[TokenKind]ParseFn{
	KindAgenda:  ParseAgenda,
	KindComment: ParseComment,
	KindKeyword: ParseKeyword,
}

// NewParser creates a new Parser object.
func NewParser(tokens []Token, parseFns map[TokenKind]ParseFn) Parser {
	return Parser{tokens: tokens, parseFns: parseFns}
}

// DefaultParser creates a new Parser object with the default parser functions.
func DefaultParser(tokens []Token) Parser {
	return NewParser(tokens, defaultParseFns)
}

type Parser struct {
	tokens   []Token
	parseFns map[TokenKind]ParseFn
}

func (p Parser) Parse() ([]Node, error) {
	_, nodes, err := p.parseMany(0)
	return nodes, err
}

// parseOne parses multiple Nodes and returns the number of tokens consumed,
// parsed Node and error. The argument i indicates the index of the token
// that will start parsing.
func (p *Parser) parseMany(i int) (int, []Node, error) {
	var (
		nodes []Node
		start = i
	)
	for i < len(p.tokens) {
		fn, ok := p.parseFns[p.tokens[i].kind]
		if !ok {
			return 0, nil, fmt.Errorf("unknown token: kind=%v", p.tokens[i].kind)
		}
		consumed, node, err := fn(p, i)
		if err != nil {
			return 0, nil, err
		}
		nodes = append(nodes, node)
		i += consumed
	}
	return i - start, nodes, nil
}
