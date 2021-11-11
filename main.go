package main

import (
	"github.com/ThreadedStream/preprocessor/preprocessor"
	"github.com/ThreadedStream/preprocessor/tokenizer"
	"os"
)

func main() {
	//var tokenizer, err = Init("arith.h")
	//if err != nil {
	//	panic(err)
	//}
	//
	//var _ = tokenizer.RetrieveMacroSymbolTable()
	//err = tokenizer.Rewind()
	//if err != nil {
	//	panic(err)
	//}
	//
	//var char = tokenizer.Next()
	//fmt.Printf("%d", char)
	//tokenizer.Finalize()

	var stream, err = os.Open("./sources/macro_flood.h")
	if err != nil {
		panic(err)
	}
	tok := tokenizer.Init(stream)
	tok.Tokenize()

	p := preprocessor.Init(tok)

	p.Preprocess()

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
