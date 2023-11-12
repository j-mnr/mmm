package token

// Type represents the supported tokens of the mmm language.
type Type uint8

// mmm supported types.
const (
	TypeIllegal Type = iota
	TypeEOF
	TypeIdent
	TypeInt
	TypeAssign
	TypePlus
	TypeMinus
	TypeStar
	TypeSlash
	TypeComma
	TypeSemicolon
	TypeLParen
	TypeRParen
	TypeLBrace
	TypeRBrace
	TypeFn
	TypeLet
	TypeBang
	TypeLT
	TypeGT
	TypeIf
	TypeElse
	TypeBool
	TypeEQ
	TypeNotEQ
	TypeReturn
	TypeString
	TypeLBrakt
	TypeRBrakt

	// TypeLookup isn't an actual type but a convenience for the [lexer.Lexer] to
	// pass in a literal value to get a correct [Token].
	TypeLookup Type = 255
)

func (t Type) String() string {
	return types[t]
}

var types = [...]string{
	"Illegal",
	"EOF",
	"Ident",
	"Int",
	"Assign",
	"Plus",
	"Minus",
	"Star",
	"Slash",
	"Comma",
	"Semicolon",
	"LParen",
	"RParen",
	"LBrace",
	"RBrace",
	"Fn",
	"Let",
	"Bang",
	"LT",
	"GT",
	"If",
	"Else",
	"Bool",
	"EQ",
	"NotEQ",
	"Return",
	"String",
	"LBrakt",
	"RBrakt",
}

// Token is one of the supported types of the mmm programming language with some
// debugging and metadata information.
type Token struct {
	typ Type
	lit string
}

func New(t Type, literal string) Token {
	if t == TypeLookup {
		switch literal {
		case "let":
			return Token{typ: TypeLet, lit: "let"}
		case "fn":
			return Token{typ: TypeFn, lit: "fn"}
		case "return":
			return Token{typ: TypeReturn, lit: "return"}
		case "if":
			return Token{typ: TypeIf, lit: "if"}
		case "else":
			return Token{typ: TypeElse, lit: "else"}
		case "true":
			return Token{typ: TypeBool, lit: "true"}
		case "false":
			return Token{typ: TypeBool, lit: "false"}
		default:
			t = TypeIdent
		}
	}
	return Token{typ: t, lit: literal}
}

func (t Token) Type() Type      { return t.typ }
func (t Token) Literal() string { return t.lit }
func (t Token) String() string { return "Type: " + t.typ.String() + "Literal: " + t.lit }
