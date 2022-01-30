package org

import (
	"fmt"
)

// parseFn is a function to parse Node and return the number of tokens consumed,
// parsed Node and error. The argument i indicates the index of the token that
// will start parsing.
type parseFn = func(p *Parser, i int) (consumed int, node Node, err error)

// parseFns expresses currently supported passers.
var parseFns = map[TokenKind]parseFn{
	KindAgenda: ParseAgenda,
}

// NewParser creates a new Parser object.
func NewParser(tokens []Token) Parser {
	return Parser{tokens: tokens}
}

type Parser struct {
	tokens []Token
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
		fn, ok := parseFns[p.tokens[i].kind]
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
