package main

import (
	"fmt"
	"github.com/ThreadedStream/preprocessor/preprocessor"
	"github.com/ThreadedStream/preprocessor/tokenizer"
	"os"
	"time"
)

const (
	gen = "./gen"
)

const (
	IN0  = "./sources/macro_flood.h"
	OUT0 = "./gen/macro_flood_prep.h"

	IN1  = "./sources/arith.h"
	OUT1 = "./gen/arith_prep.h"
)

func setupPrerequisites() {
	_, err := os.Stat(gen)
	if err != nil {
		err := os.Mkdir(gen, os.ModeDir)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	setupPrerequisites()

	start := time.Now().UnixMilli()
	var macroFloodStream, err = os.Open(IN0)
	if err != nil {
		panic(err)
	}
	tok := tokenizer.Init(macroFloodStream)
	tok.Tokenize()

	p := preprocessor.Init(tok)

	p.Preprocess(OUT0)

	arithStream, err := os.Open(IN1)
	if err != nil {
		panic(err)
	}

	tok = tokenizer.Init(arithStream)
	tok.Tokenize()

	p = preprocessor.Init(tok)

	p.Preprocess(OUT1)

	end := time.Now().UnixMilli()

	took := end - start

	fmt.Printf("Preprocessing of two files took %d milliseconds\n", took)

	return
}
