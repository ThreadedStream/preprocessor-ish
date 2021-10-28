package main

import (
	. "github.com/ThreadedStream/preprocessor/parser"
	"unicode"
)

var (
	macroSymbolTable = make(map[string]string, 0)
)

func eatWhitespaces(parser *Parser) {
	ch := parser.Peek()
	for ch != EOF && unicode.IsSpace(ch) {
		parser.Next()
		ch = parser.Peek()
	}
}

func parseDefineBlock(parser *Parser) (lhs []rune, rhs []rune) {
	for {
		var char = parser.Next()
		if unicode.IsSpace(char) {
			break
		}
		lhs = append(lhs, char)
	}

	for {
		var char = parser.Next()
		if char == '\r' || char == '\n' {
			parser.EatWhitespaces()
			break
		} else if char == '\\' {
			char = parser.Next()
			var expression = parser.ConsumeUpUntil2('\r', '\n')
			// skip terminator
			parser.Next()
			// parsing multiline define block
			multilineDefineEnded := false
			for !multilineDefineEnded {
				expression = append(expression, parser.GetLine()...)
				if expression[len(expression)-2] == '\\' && unicode.IsSpace(expression[len(expression)-1]) {
					// erase backslashes, as they solely serve as a tool to indicate ending of multiline macro
					expression = expression[0 : len(expression)-2]
					expression = append(expression, '\n')
					parser.EatWhitespaces()
					continue
				} else {
					multilineDefineEnded = true
				}
			}

			rhs = append(rhs, expression...)
			goto exit
		}

		rhs = append(rhs, char)
	}

exit:
	return
}

func determineMacro(parser *Parser) string {
	var (
		token []rune
		char  = parser.Next()
	)
	for char != EOF && !unicode.IsSpace(char) {
		token = append(token, char)
		char = parser.Next()
	}
	return string(token)
}

func advanceTillNextLine(parser *Parser) {
	char := parser.Next()
	for char != EOF && char != '\r' && char != '\n' {
		char = parser.Next()
	}

	return
}

func main() {
	var parser, err = Init("sample.h")
	if err != nil {
		panic(err)
	}

	for {
		var char = parser.Next()
		if char == EOF {
			break
		}

		switch char {
		case '#':
			macro := determineMacro(parser)
			switch macro {
			case "define":
				lhs, rhs := parseDefineBlock(parser)
				macroSymbolTable[string(lhs)] = string(rhs)
			case "include":
				advanceTillNextLine(parser)
			default:
				continue
			}
			break
		}
	}
}
