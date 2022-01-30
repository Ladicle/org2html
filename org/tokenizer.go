package org

import (
	"bufio"
	"fmt"
	"io"
)

type TokenKind string

const (
	KindAgenda TokenKind = "agenda"
)

// NewToken creates new Token object.
func NewToken(kind TokenKind, itemNumber int, matchedValues []string) Token {
	return Token{
		kind: kind,
		num:  itemNumber,
		vals: matchedValues,
	}
}

// Token is a structure to store token information.
type Token struct {
	// kind is the kind of Token.
	kind TokenKind
	// num is the number of the matched items.
	num int
	// vals is the matched values. It can contain multiple item values.
	vals []string
}

// lexFn is a lexer function which returns token and flag.
type lexFn = func(line string) (t Token, ok bool)

// lexFns expresses currently supported lexers.
var lexFns = []lexFn{
	LexAgenda,
}

// Tokenize scan each line and Tokenize them with the lexFns.
func Tokenize(in io.Reader) ([]Token, error) {
	var (
		scanner = bufio.NewScanner(in)
		tokens  []Token
	)
nextLine:
	for scanner.Scan() {
		line := scanner.Text()
		// try all lexFns
		for _, lexFn := range lexFns {
			if token, ok := lexFn(line); ok {
				tokens = append(tokens, token)
				continue nextLine
			}
		}
		return nil, fmt.Errorf("no lexers can parse %q", line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}
