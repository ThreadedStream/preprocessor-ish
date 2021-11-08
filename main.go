package main

import (
	"github.com/ThreadedStream/preprocessor/parser"
	"os"
)

func main() {
	//var parser, err = Init("sample.h")
	//if err != nil {
	//	panic(err)
	//}
	//
	//var _ = parser.RetrieveMacroSymbolTable()
	//err = parser.Rewind()
	//if err != nil {
	//	panic(err)
	//}
	//
	//var char = parser.Next()
	//fmt.Printf("%d", char)
	//parser.Finalize()

	var stream, err = os.Open("sample.h")
	if err != nil {
		panic(err)
	}
	tokenizer := parser.Init(stream)
	tokenizer.Tokenize()

	return
}

func testStreamWrite(path string) {
	var file, err = os.OpenFile(path, os.O_RDWR, 666)
	if err != nil {
		panic(err)
	}

	_, err = file.Write([]byte("hellothereb"))
	if err != nil {
		panic(err)
	}

	file.Close()
}
