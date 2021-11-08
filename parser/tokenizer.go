package parser

import (
	"fmt"
	"io"
	sc "text/scanner"
	"unicode"
)

type Tokenizer struct {
	scanner  *sc.Scanner
	stream   io.Reader
	position sc.Position
	curr     rune
	prev     rune
	next     rune
}

func Init(stream io.Reader) *Tokenizer {
	var p = &Tokenizer{
		scanner: new(sc.Scanner),
	}
	p.scanner = p.scanner.Init(stream)
	p.stream = stream
	p.curr = BOF
	p.prev = BOF
	p.next = BOF

	return p
}

func (tok *Tokenizer) Tokenize() {
	var tokens []Token
	char := tok.Next()

	for {
		if tok.curr == EOF {
			break
		}

		switch tok.curr {
		// arithmetic operators
		case '+', '-', '*', '/', '%':
			// handle assign-arithmetic operator
			if tok.next == '=' {
				tokens = append(tokens, Token{
					fmt.Sprintf("%c%c", tok.next, tok.curr),
					tok.position.Line,
					tok.position.Column,
					false,
					"",
				})
			} else {
				tokens = append(tokens, Token{
					string(tok.curr),
					tok.position.Line,
					tok.position.Column,
					false,
					"",
				})
			}
			break

		case '=':
			if tok.next == '=' {
				tokens = append(tokens, Token{
					"==",
					tok.position.Line,
					tok.position.Column,
					false,
					"",
				})
			} else {
				tokens = append(tokens, Token{
					"=",
					tok.position.Line,
					tok.position.Column,
					false,
					"",
				})
			}
			break

		case '>', '<':
			if tok.next == '<' && tok.curr == '<' {
				tokens = append(tokens, Token{
					"<<",
					tok.position.Line,
					tok.position.Column - 1,
					false,
					"",
				})
			} else if tok.next == '>' && tok.curr == '>' {
				tokens = append(tokens, Token{
					">>",
					tok.position.Line,
					tok.position.Column - 1,
					false,
					"",
				})
			} else if tok.next == '=' {
				tokens = append(tokens, Token{
					fmt.Sprintf("%c%c", char, '='),
					tok.position.Line,
					tok.position.Column - 1,
					false,
					"",
				})
			}
			break
		case '#':
			macro := tok.determineMacro()
			switch macro {
			case "include":
				line := tok.position.Line
				column := tok.position.Column - len("include") + 1
				tokens = append(tokens, Token{
					"include",
					line,
					column,
					true,
					tok.parseIncludeStmt(),
				})
			}

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			var num []rune

			num = append(num, tok.curr)

			char = tok.Next()

			for unicode.IsNumber(tok.curr) {
				num = append(num, char)
				char = tok.Next()
			}

			tokens = append(tokens, Token{
				string(num),
				tok.position.Line,
				tok.position.Column,
				false,
				"",
			})

			continue

			// parse identifier
		default:
			if unicode.IsSpace(tok.curr) {
				break
			}

			var ident []rune
			for isalphanum(tok.curr) {
				ident = append(ident, tok.curr)
				tok.curr = tok.Next()
			}
			tokens = append(tokens, Token{
				string(ident),
				tok.position.Line,
				tok.position.Column - len(ident) + 1,
				false,
				"",
			})
			continue
		}

		tok.Next()
	}
}

func (tok *Tokenizer) determineMacro() string {
	var (
		token []rune
		char  = tok.Next()
	)
	for char != EOF && !unicode.IsSpace(char) {
		token = append(token, char)
		char = tok.Next()
	}
	return string(token)
}

// parseIncludeStmt returns string representation of an include's target
// par exemple:
// #include <stdio.h>
// #include "localfile.h"

func (tok *Tokenizer) parseIncludeStmt() string {
	tok.EatWhitespaces()

	var (
		macroBody []rune
		delim     rune
	)

	if tok.curr == '<' || tok.curr == '"' {
		if tok.curr == '<' {
			delim = '>'
		} else {
			delim = '"'
		}

		for tok.curr != EOF && tok.curr != delim {
			macroBody = append(macroBody, tok.curr)
			tok.Next()
		}

		if tok.curr != delim {
			panic(fmt.Errorf("%d:%d expected delimiter %c at the end of include statement",
				tok.position.Line,
				tok.position.Column,
				delim))
		}
	} else {
		panic(fmt.Errorf("%d:%d included file must be enclosed by < or \" ", tok.position.Line, tok.position.Column))
	}

	macroBody = append(macroBody, delim)
	return string(macroBody)
}

func (tok *Tokenizer) Next() rune {
	tok.prev = tok.curr
	tok.curr = tok.scanner.Next()
	tok.next = tok.scanner.Peek()
	tok.position = tok.scanner.Pos()

	return tok.curr
}

func (tok *Tokenizer) EatWhitespaces() {
	for tok.curr != EOF && unicode.IsSpace(tok.curr) {
		tok.Next()
	}
}

func (tok *Tokenizer) Peek() rune {
	if tok.next == BOF {
		tok.curr = tok.scanner.Next()
		tok.next = tok.scanner.Peek()
	}
	return tok.next
}

func (tok *Tokenizer) Prev() rune {
	return tok.prev
}

func (tok *Tokenizer) Curr() rune {
	return tok.curr
}

func (tok *Tokenizer) Pos() sc.Position {
	return tok.position
}

func isKeyword(token string) bool {
	lo := 0
	hi := len(keywords)

	for lo <= hi {
		mid := (hi + lo) >> 1
		if keywords[mid] == token {
			return true
		} else if keywords[mid] < token {
			lo = mid + 1
			continue
		} else {
			hi = mid - 1
			continue
		}
	}

	return false
}

func isalphanum(r rune) bool {
	return (r >= '0' && r <= '9') || unicode.IsLetter(r)
}
