package preprocessor

import (
	"github.com/ThreadedStream/preprocessor/tokenizer"
)

// EOT stands for End Of Tokens
type EOT tokenizer.Token

type Preprocessor struct {
	// QUESTION(threadedstream): Do i need to store the whole tokenizer if i only need tokens?
	tokenizer *tokenizer.Tokenizer
	tOffset   int // tOffset represents an index into tokenizer.Tokens array
	prev      tokenizer.Token
	curr      tokenizer.Token
	next      tokenizer.Token
}

func Init(tokenizer *tokenizer.Tokenizer) *Preprocessor {
	if len(tokenizer.Tokens) == 0 {
		panic("you've gotta fill it with some tokens, buddy")
	}
	return &Preprocessor{
		tokenizer,
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

func (p *Preprocessor) consumeToken() tokenizer.Token {
	oldtOffset := p.tOffset
	p.tOffset++
	return p.tokenizer.Tokens[oldtOffset]
}

//
func (p *Preprocessor) peek() tokenizer.Token {
	lastTokensIdx := len(p.tokenizer.Tokens) - 1
	if p.tOffset == lastTokensIdx {
		return tokenizer.Token(EOT{})
	} else {
		return p.tokenizer.Tokens[p.tOffset]
	}
}

func (p *Preprocessor) Preprocess() {

}
