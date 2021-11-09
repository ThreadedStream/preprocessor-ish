package tokenizer

import (
	"fmt"
	"io"
	sc "text/scanner"
	"unicode"
)

// BOF stands for "Beginning of file"
const (
	EOF = -1
	BOF = -2
)

type Tokenizer struct {
	scanner  *sc.Scanner // scanner is a stdlib implementation of a tokenizer
	stream   io.Reader   // stream is an actual source file
	position sc.Position // position indicates a current position in the stream
	Tokens   []Token     // tokens represents a slice of tokens consumed so far
	prev     rune        // prev represents a rune prior to curr
	curr     rune        // curr represents a rune currently being under consideration
	next     rune        // next represents a next-to-curr rune
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
	tok.Next()

	for {
		if tok.curr == EOF {
			break
		}

		switch tok.curr {
		// arithmetic operators, comments
		case '+', '-', '*', '/', '%':
			if tok.next == '=' {
				tok.Tokens = append(tok.Tokens, Token{
					fmt.Sprintf("%c%c", tok.next, tok.curr),
					tok.position.Line,
					tok.position.Column,
					false,
					false,
					"",
				})
				tok.Next()
			} else if tok.curr == '/' && tok.next == '/' {
				tok.EatSingleLineComment()
				continue
			} else if tok.curr == '/' && tok.next == '*' {
				tok.EatMultiLineComment()
				continue
			} else {
				tok.Tokens = append(tok.Tokens, Token{
					string(tok.curr),
					tok.position.Line,
					tok.position.Column - 1,
					false,
					false,
					"",
				})
			}
			break
		// equality, assignment operators
		case '=':
			if tok.next == '=' {
				tok.Tokens = append(tok.Tokens, Token{
					"==",
					tok.position.Line,
					tok.position.Column - 1,
					false,
					false,
					"",
				})
				tok.Next()
			} else {
				tok.Tokens = append(tok.Tokens, Token{
					"=",
					tok.position.Line,
					tok.position.Column - 1,
					false,
					false,
					"",
				})
			}
			break

			// shift operators, greater-or-equal, less-or-equal operators
			// TODO(threadedstream): add support for a recognition of greater and less than operators
		case '>', '<':
			if tok.next == '<' && tok.curr == '<' {
				tok.Tokens = append(tok.Tokens, Token{
					"<<",
					tok.position.Line,
					tok.position.Column - 1,
					false,
					false,
					"",
				})
				tok.Next()
			} else if tok.next == '>' && tok.curr == '>' {
				tok.Tokens = append(tok.Tokens, Token{
					">>",
					tok.position.Line,
					tok.position.Column - 1,
					false,
					false,
					"",
				})
				tok.Next()
			} else if tok.next == '=' {
				tok.Tokens = append(tok.Tokens, Token{
					fmt.Sprintf("%c%c", tok.curr, '='),
					tok.position.Line,
					tok.position.Column - 1,
					false,
					false,
					"",
				})
				tok.Next()
			}
			break

		// preprocessor directives
		case '#':
			line := tok.position.Line
			column := tok.position.Column
			macro := tok.determineMacro()
			switch macro {
			case "include":
				tok.Tokens = append(tok.Tokens, Token{
					"include",
					line,
					column,
					true,
					false,
					tok.getIncludeBody(),
				})
			case "define":
				lhs, rhs := tok.getDefineBody()
				tok.Tokens = append(tok.Tokens, Token{
					"define",
					line,
					column,
					true,
					false,
					map[string]string{
						string(lhs): string(rhs),
					},
				})
			}
		// numbers
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			var num []rune

			num = append(num, tok.curr)

			tok.Next()

			for unicode.IsNumber(tok.curr) {
				num = append(num, tok.curr)
				tok.Next()
			}

			tok.Tokens = append(tok.Tokens, Token{
				string(num),
				tok.position.Line,
				tok.position.Column,
				false,
				false,
				"",
			})

			continue

		// parse identifier
		default:
			if unicode.IsSpace(tok.curr) {
				break
			}

			var (
				ident  []rune
				line   = tok.position.Line
				column = tok.position.Column - 1
			)
			if isAllowedIdentRune(tok.curr) {
				for isAllowedIdentRune(tok.curr) {
					ident = append(ident, tok.curr)
					tok.Next()
				}
			}
			if len(ident) == 0 {
				tok.Tokens = append(tok.Tokens, Token{
					string(tok.curr),
					line,
					column,
					false,
					false,
					"",
				})
				tok.Next()
			} else {
				tok.Tokens = append(tok.Tokens, Token{
					string(ident),
					line,
					column,
					false,
					isKeyword(string(ident)),
					"",
				})
			}

			continue
		}

		tok.Next()
	}
}

func (tok *Tokenizer) EatSingleLineComment() {
	for tok.curr != EOF && tok.curr != '\r' && tok.curr != '\n' {
		tok.Next()
	}
	tok.EatWhitespaces()
}

func (tok *Tokenizer) EatMultiLineComment() {
	for tok.curr != EOF && !(tok.curr == '*' && tok.next == '/') {
		tok.Next()
	}
	// TODO(threadedstream): probably, should be replaced with call to some
	// (not yet created) SkipChars(), as that would be much cleaner
	tok.Next()
	tok.Next()
	tok.EatWhitespaces()
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

func (tok *Tokenizer) ConsumeUpUntil2(term0, term1 rune) (token []rune) {
	for tok.curr != term0 && tok.curr != term1 {
		token = append(token, tok.curr)
		tok.Next()
	}

	return
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
	hi := len(keywords) - 1

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

func isAllowedIdentRune(r rune) bool {
	return isalphanum(r) || r == '_'
}

func isalphanum(r rune) bool {
	return (r >= '0' && r <= '9') || unicode.IsLetter(r)
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

func (tok *Tokenizer) getDefineBody() (lhs []rune, rhs []rune) {
	for {
		tok.Next()
		if unicode.IsSpace(tok.curr) {
			break
		}
		lhs = append(lhs, tok.curr)
	}

	// skip all whitespaces
	tok.EatWhitespaces()

	for {
		switch tok.curr {
		case '\\':
			// skip \ and whitespaces
			tok.Next()
			tok.EatWhitespaces()
			expr := tok.getMultilineDefineBody()
			if len(expr) == 0 {
				goto exit
			}
			rhs = append(rhs, expr...)
			goto exit
		case '\n', '\r':
			// we're done with a one-liner
			goto exit
		default:
			rhs = append(rhs, tok.curr)
			tok.Next()
			continue
		}
	}

exit:
	return

}

func (tok *Tokenizer) getMultilineDefineBody() (expr []rune) {
	// postpone coalescing of consumedSoFar into expr

	for {
		expr = append(expr, tok.ConsumeUpUntil2('\\', '\n')...)
		// dummy define expression
		if len(expr) == 0 {
			// no error throw
			break
		}
		// consume last character residing at the end of a line
		expr = append(expr, tok.curr)
		// characters \n and \r at the end of expression serve as
		// indicators of define expression's completion
		if expr[len(expr)-1] == '\n' || expr[len(expr)-1] == '\r' {
			return
		}

		// rip the last character off of the expression
		expr = expr[0 : len(expr)-1]
		tok.Next()
		tok.EatWhitespaces()
	}

	return
}

// parseIncludeStmt returns string representation of an include's target
// par exemple:
// #include <stdio.h>
// #include "localfile.h"

func (tok *Tokenizer) getIncludeBody() string {
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
