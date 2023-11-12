package lexer

import (
	"mmm/token"
)

// Lexer takes an input and turns it into a slice of [token.Token]. This
// implementation doesn't care about whitespace.
type Lexer struct {
	// input is the actual value we are trying to tokenize.
	input string
	// cPos is the current cPos in the input string.
	cPos uint
	// nPos is the next position in the input string.
	nPos uint
	// ch is the char at cPos.
	ch byte
}

// New returns a Lexer that will parse the input token by token.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// NextToken provides the next token in the Lexer's input.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.eatWhitespace()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			tok = token.New(token.TypeEQ, "==")
			l.readChar()
			break
		}
		tok = token.New(token.TypeAssign, string(l.ch))
	case '+':
		tok = token.New(token.TypePlus, string(l.ch))
	case '-':
		tok = token.New(token.TypeMinus, string(l.ch))
	case '*':
		tok = token.New(token.TypeStar, string(l.ch))
	case '/':
		tok = token.New(token.TypeSlash, string(l.ch))
	case '!':
		if l.peekChar() == '=' {
			tok = token.New(token.TypeNotEQ, "!=")
			l.readChar()
			break
		}
		tok = token.New(token.TypeBang, string(l.ch))
	case '<':
		tok = token.New(token.TypeLT, string(l.ch))
	case '>':
		tok = token.New(token.TypeGT, string(l.ch))
	case ';':
		tok = token.New(token.TypeSemicolon, string(l.ch))
	case '(':
		tok = token.New(token.TypeLParen, string(l.ch))
	case ')':
		tok = token.New(token.TypeRParen, string(l.ch))
	case ',':
		tok = token.New(token.TypeComma, string(l.ch))
	case '{':
		tok = token.New(token.TypeLBrace, string(l.ch))
	case '}':
		tok = token.New(token.TypeRBrace, string(l.ch))
	case '[':
		tok = token.New(token.TypeLBrakt, string(l.ch))
	case ']':
		tok = token.New(token.TypeRBrakt, string(l.ch))
	case 0:
		tok = token.New(token.TypeEOF, "")
	case '"':
		tok = token.New(token.TypeString, l.readString())
	default:
		switch {
		case isLetter(l.ch):
			return token.New(token.TypeLookup, l.readIdentifier())
		case isDigit(l.ch):
			return token.New(token.TypeInt, l.readNumber())
		default:
			tok = token.New(token.TypeIllegal, string(l.ch))
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	l.ch = 0
	if l.nPos < uint(len(l.input)) {
		l.ch = l.input[l.nPos]
	}
	l.cPos = l.nPos
	l.nPos++
}

func (l Lexer) peekChar() byte {
	if l.nPos == uint(len(l.input)) {
		return 0
	}
	return l.input[l.nPos]
}

func (l *Lexer) readIdentifier() string {
	position := l.cPos
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.cPos]
}

func (l *Lexer) readNumber() string {
	start := l.cPos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.cPos]
}

func (l *Lexer) eatWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(b byte) bool { return '0' <= b && b <= '9' }

func (l *Lexer) readString() string {
	l.readChar()
	start := l.cPos
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	return l.input[start:l.cPos]
}
