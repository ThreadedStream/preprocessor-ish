package parser

import (
	"io"
	"strings"
)

const (
	SRC0 = `
		#include <stdlib.h>
		#include <stdio.h>
	
	
		int main(int argc, const char* argv[]) {
			int a = 32;
			int b = 122;
		
			int c = a + b;
			c += b;
			c /= 3;
			c %= 2;
			c = a % 5;
		
			c = b / 2;
		
		}
`
)

func SrcToReader(stringSrc string) io.Reader {
	return strings.NewReader(stringSrc)
}
