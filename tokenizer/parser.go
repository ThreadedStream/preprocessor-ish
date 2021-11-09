package tokenizer

import (
	"os"
	sc "text/scanner"
	"unicode"
)

type Parser struct {
	scanner *sc.Scanner
	stream  *os.File
	curr    rune
	prev    rune
	next    rune
}

func (p *Parser) Rewind() error {
	_, err := p.stream.Seek(0, 0)
	if err != nil {
		return err
	}

	// reinitialize a scanner
	p.scanner = p.scanner.Init(p.stream)
	return nil
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

func (p *Parser) EatSingleLineComment() {
	for p.curr != EOF && p.curr != '\r' && p.curr != '\n' {
		p.Next()
	}
	p.EatWhitespaces()
}

func (p *Parser) EatMultiLineComment() {
	for p.curr != EOF && !(p.curr == '*' && p.Peek() == '/') {
		p.Next()
	}
	// TODO(threadedstream): probably, should be replaced with call to some
	// (not yet created) SkipChars(), as that would be much cleaner
	p.Next()
	p.Next()
	p.EatWhitespaces()
}

// RetrieveMacroSymbolTable will analyze not-yet-preprocessed
// source code and fill "macro symbol table" with mappings
// from symbol names to their respective values the former represent
func (p *Parser) RetrieveMacroSymbolTable() map[string]string {
	var newSrc []rune

	macroSymbolTable := make(map[string]string, 0)

	var (
		ident []rune
		char  = p.Next()
	)

	for {
		char = p.Curr()
		if char == EOF {
			break
		}

		switch char {
		case '#':
			macro := p.determineMacro()
			switch macro {
			case "define":
				lhs, rhs := p.parseDefineExpression()
				macroSymbolTable[string(lhs)] = string(rhs)
				p.EatWhitespaces()
				continue
			case "include":
				newSrc = append(newSrc, p.advanceTillNextLine()...)
				continue
			default:
				continue
			}
		case '/':
			// here comes a comment -- just skip it
			switch p.Peek() {
			case '/':
				// just a one-liner
				p.EatSingleLineComment()
				continue
			case '*':
				p.Next()
				// here we're dealing with a multi-line monster
				p.EatMultiLineComment()
				continue
			}
		default:
			ident = append(ident, char)
			newSrc = append(newSrc, char)
		}
	}

	return macroSymbolTable
}

func replaceInFile(old, new, src []rune, currPos int) {
	var (
		oldLen = len(old)
		newLen = len(new)
		offset = currPos - oldLen
	)

	// the most desired case
	if oldLen == newLen {
		newIdx := 0
		for offset < currPos {
			src[offset] = new[newIdx]
			offset++
			newIdx++
		}
	} else if oldLen > newLen {
		newIdx := 0
		for offset < newLen {
			src[offset] = new[newIdx]
			offset++
			newIdx++
		}
	} else {
		newIdx := 0
		for offset < oldLen {
			src[offset] = new[newIdx]
			offset++
			newIdx++
		}
		//newOldDiff := newLen - oldLen
	}

}

func (p *Parser) Preprocess(macroSymbolTable map[string]string) {
	err := p.Rewind()
	if err != nil {
		panic(err)
	}

	var char = p.Next()
	for {
		char = p.Curr()
		if char == EOF {
			break
		}

		switch char {
		case '#':
			macro := p.determineMacro()
			switch macro {
			case "define":
				// ignore define macros
				p.advanceTillNextLine()
			case "include":
				p.advanceTillNextLine()
				p.EatWhitespaces()
				continue
			default:
				continue
			}
		case '/':
			// here comes a comment -- just skip it
			switch p.Peek() {
			case '/':
				// just a one-liner
				p.EatSingleLineComment()
				continue
			case '*':
				p.Next()
				// here we're dealing with a multi-line monster
				p.EatMultiLineComment()
				continue
			}
		default:
			p.Next()
			continue
		}
	}

}

func (p *Parser) Finalize() {
	// for now, finalize only will close a file handle
	p.stream.Close()
}

// types of expression parseDefineBlock is able to recognize
/*
	#define IDENT VAL
	#define IDENT \
		expr \
		expr cont.\
		expr end [\] ([\] means optionality of \ character)

	#define IDENT expr \
	 	expr cont. \
		more expr  \
		expr end [\]
*/
func (p *Parser) parseDefineExpression() (lhs []rune, rhs []rune) {
	for {
		var char = p.Next()
		if unicode.IsSpace(char) {
			break
		}
		lhs = append(lhs, char)
	}

	// skip all whitespaces
	p.EatWhitespaces()

	for {
		char := p.Curr()

		switch char {
		case '\\':
			// skip \ and whitespaces
			p.Next()
			p.EatWhitespaces()
			expr := p.parseMultilineDefineExpression(rhs)
			if len(expr) == 0 {
				goto exit
			}
			rhs = append(rhs, expr...)
			goto exit
		case '\n', '\r':
			// we're done with a one-liner
			goto exit
		default:
			rhs = append(rhs, char)
			p.Next()
			continue
		}
	}

exit:
	return
}

func (p *Parser) parseMultilineDefineExpression(consumedSoFar []rune) (expr []rune) {
	// postpone coalescing of consumedSoFar into expr

	for {
		expr = append(expr, p.ConsumeUpUntil2('\\', '\n')...)
		// dummy define expression
		if len(expr) == 0 {
			// no error throw
			break
		}
		// consume last character residing at the end of a line
		expr = append(expr, p.Curr())
		// characters \n and \r at the end of expression serve as
		// indicators of define expression's completion
		if expr[len(expr)-1] == '\n' || expr[len(expr)-1] == '\r' {
			return
		}

		// rip the last character off of the expression
		expr = expr[0 : len(expr)-1]
		p.Next()
		p.EatWhitespaces()
	}

	return
}

func (p *Parser) determineMacro() string {
	var (
		token []rune
		char  = p.Next()
	)
	for char != EOF && !unicode.IsSpace(char) {
		token = append(token, char)
		char = p.Next()
	}
	return string(token)
}

func (p *Parser) advanceTillNextLine() []rune {
	var buf []rune

	char := p.Next()
	buf = append(buf, char)
	for char != EOF && char != '\r' && char != '\n' {
		char = p.Next()
		buf = append(buf, char)
	}

	return buf
}
