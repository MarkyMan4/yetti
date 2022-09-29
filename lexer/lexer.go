package lexer

import (
	"unicode"

	"github.com/MarkyMan4/yetti/token"
)

type Lexer struct {
	curPos  int
	readPos int
	curChar rune
	chars   []rune
}

func NewLexer(input string) *Lexer {
	l := &Lexer{chars: []rune(input)}
	l.nextChar()
	return l
}

func (l *Lexer) nextChar() {
	if l.readPos >= len(l.chars) {
		l.curChar = rune(0)
	} else {
		l.curChar = l.chars[l.readPos]
	}

	l.curPos = l.readPos
	l.readPos++
}

// peek ahead one character without increasing curPos or readPos
func (l *Lexer) peek() rune {
	var peekedChar rune

	if l.readPos >= len(l.chars) {
		peekedChar = rune(0)
	} else {
		peekedChar = l.chars[l.readPos]
	}

	return peekedChar
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	switch l.curChar {
	case '+':
		if l.peek() == '=' {
			tok = token.Token{Type: token.PLUSEQ, Literal: "+="}
			l.nextChar()
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.curChar)}
		}
	case '-':
		if l.peek() == '=' {
			tok = token.Token{Type: token.MINEQ, Literal: "-="}
			l.nextChar()
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.curChar)}
		}
	case '*':
		if l.peek() == '=' {
			tok = token.Token{Type: token.MULTEQ, Literal: "*="}
			l.nextChar()
		} else {
			tok = token.Token{Type: token.MULT, Literal: string(l.curChar)}
		}
	case '/':
		if l.peek() == '/' {
			l.readToEndOfLine()
		} else if l.peek() == '=' {
			tok = token.Token{Type: token.DIVEQ, Literal: "/="}
			l.nextChar()
		} else {
			tok = token.Token{Type: token.DIVIDE, Literal: string(l.curChar)}
		}
	case '=':
		if l.peek() == '=' {
			tok = token.Token{Type: token.EQ, Literal: "=="}
			l.nextChar()
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.curChar)}
		}
	case '<':
		if l.peek() == '=' {
			tok = token.Token{Type: token.LTE, Literal: "<="}
			l.nextChar()
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.curChar)}
		}
	case '>':
		if l.peek() == '=' {
			tok = token.Token{Type: token.GTE, Literal: ">="}
			l.nextChar()
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.curChar)}
		}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.curChar)}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.curChar)}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.curChar)}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.curChar)}
	case '[':
		tok = token.Token{Type: token.LBRACK, Literal: string(l.curChar)}
	case ']':
		tok = token.Token{Type: token.RBRACK, Literal: string(l.curChar)}
	case ';':
		tok = token.Token{Type: token.SEMI, Literal: string(l.curChar)}
	case ',':
		tok = token.Token{Type: token.COM, Literal: string(l.curChar)}
	case '"':
		tok = l.readString()
	case '.':
		tok = token.Token{Type: token.DOT, Literal: string(l.curChar)}
	case rune(0):
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		// read a number or an identifier
		if unicode.IsDigit(l.curChar) {
			tok = l.readIntOrFloat()
		} else {
			tok = l.readIdent()
		}
	}

	l.nextChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.curChar == ' ' || l.curChar == '\n' || l.curChar == '\t' || l.curChar == '\r' {
		l.nextChar()
	}
}

// used for comments
func (l *Lexer) readToEndOfLine() {
	for l.curChar != '\n' {
		l.nextChar()
	}
}

func (l *Lexer) readIntOrFloat() token.Token {
	literal := l.readNumber()
	var tokType string

	// if a decimal comes after the number, read the decimal and then all the numbers after it - it's a float
	if l.peek() == '.' {
		l.nextChar()
		literal += string(l.curChar)
		l.nextChar() // read past decimal before reading digits after decimal
		literal += l.readNumber()
		tokType = token.FLOAT
	} else {
		tokType = token.INT
	}

	return token.Token{Type: tokType, Literal: literal}
}

func (l *Lexer) readNumber() string {
	numTok := string(l.curChar)

	for unicode.IsDigit(l.peek()) {
		l.nextChar()
		numTok += string(l.curChar)
	}

	return numTok
}

func (l *Lexer) readString() token.Token {
	strTok := ""

	// current token is ", read past this and then read until another " is seen
	l.nextChar()

	for l.curChar != '"' && l.curChar != rune(0) {
		strTok += string(l.curChar)
		l.nextChar()
	}

	return token.Token{Type: token.STRING, Literal: strTok}
}

func (l *Lexer) readIdent() token.Token {
	literal := string(l.curChar)

	for unicode.IsLetter(l.peek()) || unicode.IsDigit(l.peek()) {
		l.nextChar()
		literal += string(l.curChar)
	}

	tokType := token.GetIdentOrKeyword(literal)

	if literal == "true" || literal == "false" {
		tokType = token.BOOLEAN
	}

	return token.Token{Type: tokType, Literal: literal}
}
