package token

type Token struct {
	Type    string
	Literal string
}

const (
	VAR     = "VAR"
	FOR     = "FOR"
	WHILE   = "WHILE"
	IF      = "IF"
	ELSE    = "ELSE"
	FUN     = "FUN"
	RETURN  = "RETURN"
	PLUS    = "+"
	MINUS   = "-"
	MULT    = "*"
	DIVIDE  = "/"
	LT      = "<"
	LTE     = "<="
	EQ      = "=="
	GT      = ">"
	GTE     = ">="
	ASSIGN  = "="
	PLUSEQ  = "+="
	MINEQ   = "-="
	MULTEQ  = "*="
	DIVEQ   = "/="
	LPAREN  = "("
	RPAREN  = ")"
	LBRACE  = "{"
	RBRACE  = "}"
	LBRACK  = "["
	RBRACK  = "]"
	SEMI    = ";"
	COM     = ","
	DQUOTE  = "\""
	DOT     = "."
	IDENT   = "IDENT"
	INT     = "INT"
	FLOAT   = "FLOAT"
	STRING  = "STRING"
	BOOLEAN = "BOOLEAN"
	EOF     = "EOF"
)

var keywords = map[string]string{
	"var":    VAR,
	"for":    FOR,
	"while":  WHILE,
	"if":     IF,
	"else":   ELSE,
	"fun":    FUN,
	"return": RETURN,
}

// lookup a value from the input and determine if it is a keyword or an identifier
func GetIdentOrKeyword(literal string) string {
	if tok, ok := keywords[literal]; ok {
		return tok
	}

	return IDENT
}
