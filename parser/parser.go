package parser

import (
	"os"
	sc "text/scanner"
)


// BOF stands for "Beginning of file"
const (
	EOF = -1
	BOF = -2
)


type Parser struct {
	scanner *sc.Scanner
	curr rune
	prev rune
	next rune
}

func Init(filename string) (*Parser,error){
	var stream, err = os.Open(filename); if err != nil {
		return nil, err
	}

	var p = &Parser{}
	p.scanner = p.scanner.Init(stream)
	p.curr = BOF
	p.prev = BOF
	p.next = BOF

	return p, nil
}

func (p *Parser) Next() {
	p.prev = p.curr
	p.curr = p.scanner.Next()
	p.next = p.scanner.Peek()
}

func (p* Parser) Peek() rune {
	if p.next == BOF{
		p.curr = p.scanner.Next()
		p.next = p.scanner.Peek()
	}
	return p.next
}

func (p *Parser) Prev() rune {
	return p.prev
}

