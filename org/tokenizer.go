package org

import (
	"bufio"
	"fmt"
	"io"
)

type TokenKind string

const (
	KindAgenda   TokenKind = "agenda"
	KindKeyword            = "keyword"
	KindComment            = "comment"
	KindHeadline           = "headline"
	KindSection            = "section"
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

// LexFn is a lexer function which returns token and flag.
type LexFn = func(line string) (t Token, ok bool)

// defaultLexFns expresses currently supported lexers.
var defaultLexFns = []LexFn{
	LexHeadline, // * <keyword> <priority> <title> <tags>
	LexKeyword,  // #+<keyword>: <val>
	LexComment,  // # <comment>
	LexAgenda,   // <agenda>: <date>

	LexSection, // *
}

// NewTokenizer creates a new Tokenizer object.
func NewTokenizer(lexFns []LexFn) Tokenizer {
	return Tokenizer{lexFns: lexFns}
}

// DefaultTokenizer creates a new Tokenizer object which has default LexFns.
func DefaultTokenizer() Tokenizer {
	return Tokenizer{lexFns: defaultLexFns}
}

type Tokenizer struct {
	lexFns []LexFn
}

// Tokenize scan each line and Tokenize them with the lexFns.
func (t Tokenizer) Tokenize(in io.Reader) ([]Token, error) {
	var (
		scanner = bufio.NewScanner(in)
		tokens  []Token
	)
nextLine:
	for scanner.Scan() {
		line := scanner.Text()
		// try all lexFns
		for _, lexFn := range t.lexFns {
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
