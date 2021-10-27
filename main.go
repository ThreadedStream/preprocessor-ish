package main

//// #include <stdio.h>
//import "C"

import (
	"os"
	sc "text/scanner"
	"unicode"
)

var (
	macroSymbolTable = make(map[string]string, 0)
)


//func parseMultilineDefineBlock(scanner *sc.Scanner) (lhs []rune, rhs []rune) {
//
//}

func eatWhitespaces(scanner *sc.Scanner) {
	ch := scanner.Peek()
	for ch != EOF && unicode.IsSpace(ch){
		scanner.Next()
		ch = scanner.Peek()
	}
}

func parseDefineBlock(scanner *sc.Scanner) (lhs []rune, rhs []rune) {

	for {
		var char = scanner.Next()
		if unicode.IsSpace(char) {
			break
		}
		lhs = append(lhs, char)
	}

	for {
		var char = scanner.Next()
		if unicode.IsSpace(char) {
			break
		} else if char == '\\' {
			eatWhitespaces(scanner)
			char = scanner.Next()
			// parsing multiline define block
			for char != '\\' && (scanner.Peek() != '\n' || scanner.Peek() != '\r'){
				rhs = append(rhs, char)
				char = scanner.Next()
			}

			break
		}

		rhs = append(rhs, char)
	}

	return
}

func determineMacro(scanner *sc.Scanner) string{
	char := scanner.Next()
	var token []rune
	for char != EOF && !unicode.IsSpace(char) {
		token = append(token, char)
		char = scanner.Next()
	}
	return string(token)
}

func advanceTillNextLine(scanner *sc.Scanner) {
	char := scanner.Next()
	for char != EOF && char != '\r' && char != '\n'{
		char = scanner.Next()
	}

	return
}

const (
	EOF = -1
)

func main(){
	var stream, err = os.Open("sample.h"); if err != nil {
		panic(err)
	}

	var scanner = &sc.Scanner{}

	scanner = scanner.Init(stream)


	for {
		var char = scanner.Next()
		if char == EOF{
			break
		}
		switch char{
		case '#':
			macro := determineMacro(scanner)
			switch macro {
			case "define":
				lhs, rhs := parseDefineBlock(scanner)
				macroSymbolTable[string(lhs)] = string(rhs)
			case "include":
				advanceTillNextLine(scanner)
			default:
				continue
			}
			break
		}
	}

}

