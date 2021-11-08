package parser

var (
	keywords = []string{
		"auto", "break", "case", "char",
		"const", "continue", "default", "define", "do",
		"double", "else", "enum", "extern",
		"float", "for", "goto", "if", "ifdef",
		"include", "int", "long", "register",
		"return", "short", "signed", "sizeof",
		"static", "struct", "switch", "typedef",
		"union", "unsigned", "void", "volatile", "while",
	}
)

type Token struct {
	value     string
	line      int
	column    int
	isMacro   bool
	MacroBody string
}
