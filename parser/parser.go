package parser

import (
	"os"
	sc "text/scanner"
	"unicode"
)

// BOF stands for "Beginning of file"
const (
	EOF = -1
	BOF = -2
)

type Parser struct {
	scanner *sc.Scanner
	curr    rune
	prev    rune
	next    rune
}

func Init(filename string) (*Parser, error) {
	var stream, err = os.Open(filename)
	if err != nil {
		return nil, err
	}

	var p = &Parser{
		scanner: new(sc.Scanner),
	}
	p.scanner = p.scanner.Init(stream)
	p.curr = BOF
	p.prev = BOF
	p.next = BOF

	return p, nil
}

func (p *Parser) Next() rune {
	p.prev = p.curr
	p.curr = p.scanner.Next()
	p.next = p.scanner.Peek()

	return p.curr
}

func (p *Parser) Peek() rune {
	if p.next == BOF {
		p.curr = p.scanner.Next()
		p.next = p.scanner.Peek()
	}
	return p.next
}

func (p *Parser) Prev() rune {
	return p.prev
}

func (p *Parser) Curr() rune {
	return p.curr
}

func (p *Parser) GetLine() (line []rune) {
	for p.curr != '\r' && p.curr != '\n' {
		line = append(line, p.curr)
		p.Next()
	}

	// include a newline
	line = append(line, p.curr)

	return
}

func (p *Parser) ConsumeUpUntil2(term0, term1 rune) (token []rune) {
	for p.curr != term0 && p.curr != term1 {
		token = append(token, p.curr)
		p.Next()
	}

	return
}

func (p *Parser) EatWhitespaces() {
	for p.curr != EOF && unicode.IsSpace(p.curr) {
		p.Next()
	}
}
