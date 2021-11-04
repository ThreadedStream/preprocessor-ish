package main

import (
	"fmt"
	. "github.com/ThreadedStream/preprocessor/parser"
)

func main() {
	var parser, err = Init("sample.h")
	if err != nil {
		panic(err)
	}

	var _ = parser.RetrieveMacroSymbolTable()
	err = parser.Rewind()
	if err != nil {
		panic(err)
	}

	var char = parser.Next()
	fmt.Printf("%d", char)
	parser.Finalize()
	return
}
