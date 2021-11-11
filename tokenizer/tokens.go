package tokenizer

const (
	IDENT = iota
	NUMBER
	KEYWORD
	OPERATOR
	MACRO
	GARBAGE
)

type TokenType int

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

// Token TODO(threadedstream): get rid of IsMacro field
type Token struct {
	Text      string
	TokType   TokenType
	Line      int
	Column    int
	IsMacro   bool
	IsKeyword bool
	MacroBody interface{}
}
