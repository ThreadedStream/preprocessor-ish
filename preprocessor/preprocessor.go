package preprocessor

import (
	"github.com/ThreadedStream/preprocessor/tokenizer"
	"os"
)

// EOT stands for End Of Tokens
type EOT tokenizer.Token

type Preprocessor struct {
	// QUESTION(threadedstream): Do i need to store the whole tokenizer if i only need tokens?
	tokenizer      *tokenizer.Tokenizer
	macroSymbTable map[string]string
	tOffset        int // tOffset represents an index into tokenizer.Tokens array
	prev           tokenizer.Token
	curr           tokenizer.Token
	next           tokenizer.Token
}

func Init(tokenizer *tokenizer.Tokenizer) *Preprocessor {
	if len(tokenizer.Tokens) == 0 {
		panic("you've gotta fill it with some tokens, buddy")
	}
	return &Preprocessor{
		tokenizer,
		make(map[string]string, 0),
		0,
		tokenizer.Tokens[0],
		tokenizer.Tokens[0],
		tokenizer.Tokens[0],
	}
}

func (p *Preprocessor) Next() tokenizer.Token {
	p.prev = p.curr
	p.curr = p.consumeToken()
	p.next = p.peek()

	return p.curr
}

func (p *Preprocessor) Preprocess() {
	var source []rune

	p.gatherMacroSymbols()

	p.Next()

	for !p.isEOT(p.curr) {
		switch p.curr.TokType {
		case tokenizer.IDENT:
			if val, ok := p.macroSymbTable[p.curr.Text]; ok {
				source = append(source, []rune(val)...)
				source = append(source, ' ')
			} else {
				source = append(source, []rune(p.curr.Text)...)
				source = append(source, ' ')
			}
		case tokenizer.MACRO:
			if p.curr.Text == "define" {
				break
			} else if p.curr.Text == "include" {
				source = append(source, '#')
				source = append(source, []rune(p.curr.Text)...)
				val := p.curr.MacroBody.(string)
				source = append(source, ' ')
				source = append(source, []rune(val)...)
			}
		default:
			source = append(source, []rune(p.curr.Text)...)
			source = append(source, ' ')
		}

		p.Next()
	}

	outstream, err := os.OpenFile("source/macro_flood_test.h", os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}

	outstream.Write([]byte(string(source)))
}

func (p *Preprocessor) gatherMacroSymbols() {
	for _, val := range p.tokenizer.Tokens {
		if val.TokType == tokenizer.MACRO {
			if parameterPack, ok := val.MacroBody.([]string); ok {
				p.macroSymbTable[parameterPack[0]] = parameterPack[1]
			} else {
				continue
			}
		}
	}
}

func (p *Preprocessor) consumeToken() tokenizer.Token {
	oldtOffset := p.tOffset
	p.tOffset++
	return p.tokenizer.Tokens[oldtOffset]
}

func (p *Preprocessor) isEOT(token tokenizer.Token) bool {
	return token.Line == -1 && token.Column == -1
}

func (p *Preprocessor) peek() tokenizer.Token {
	lastTokensIdx := len(p.tokenizer.Tokens) - 1
	if p.tOffset == lastTokensIdx {
		return tokenizer.Token(EOT{})
	} else {
		return p.tokenizer.Tokens[p.tOffset+1]
	}
}
