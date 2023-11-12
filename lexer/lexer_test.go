package lexer_test

import (
	"mmm/is"
	"mmm/lexer"
	"mmm/token"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		input string
		toks  []token.Token
	}{
		"Let statement": {
			input: `
let five = 5;
let ten = 10;
let add = fn(x, y) {
	return x + y;
};
let result = add(five, ten);`,
			toks: []token.Token{
				token.New(token.TypeLet, "let"),
				token.New(token.TypeIdent, "five"),
				token.New(token.TypeAssign, "="),
				token.New(token.TypeInt, "5"),
				token.New(token.TypeSemicolon, ";"),
				token.New(token.TypeLet, "let"),
				token.New(token.TypeIdent, "ten"),
				token.New(token.TypeAssign, "="),
				token.New(token.TypeInt, "10"),
				token.New(token.TypeSemicolon, ";"),
				token.New(token.TypeLet, "let"),
				token.New(token.TypeIdent, "add"),
				token.New(token.TypeAssign, "="),
				token.New(token.TypeFn, "fn"),
				token.New(token.TypeLParen, "("),
				token.New(token.TypeIdent, "x"),
				token.New(token.TypeComma, ","),
				token.New(token.TypeIdent, "y"),
				token.New(token.TypeRParen, ")"),
				token.New(token.TypeLBrace, "{"),
				token.New(token.TypeReturn, "return"),
				token.New(token.TypeIdent, "x"),
				token.New(token.TypePlus, "+"),
				token.New(token.TypeIdent, "y"),
				token.New(token.TypeSemicolon, ";"),
				token.New(token.TypeRBrace, "}"),
				token.New(token.TypeSemicolon, ";"),
				token.New(token.TypeLet, "let"),
				token.New(token.TypeIdent, "result"),
				token.New(token.TypeAssign, "="),
				token.New(token.TypeIdent, "add"),
				token.New(token.TypeLParen, "("),
				token.New(token.TypeIdent, "five"),
				token.New(token.TypeComma, ","),
				token.New(token.TypeIdent, "ten"),
				token.New(token.TypeRParen, ")"),
				token.New(token.TypeSemicolon, ";"),
				token.New(token.TypeEOF, ""),
			},
		},
		"Invalid parsed code gives valid tokens": {
			input: `!-/*5;`,
			toks: []token.Token{
				token.New(token.TypeBang, "!"),
				token.New(token.TypeMinus, "-"),
				token.New(token.TypeSlash, "/"),
				token.New(token.TypeStar, "*"),
				token.New(token.TypeInt, "5"),
				token.New(token.TypeSemicolon, ";"),
			},
		},
		"Invalid comparisons gives valid tokens": {
			input: `5 < 10 > 5;`,
			toks: []token.Token{
				token.New(token.TypeInt, "5"),
				token.New(token.TypeLT, "<"),
				token.New(token.TypeInt, "10"),
				token.New(token.TypeGT, ">"),
				token.New(token.TypeInt, "5"),
			},
		},
		"If/else statement with bools": {
			input: `
if (5 < 10) {
	return true;
} else {
	return false;
}`,
			toks: []token.Token{
				token.New(token.TypeIf, "if"),
				token.New(token.TypeLParen, "("),
				token.New(token.TypeInt, "5"),
				token.New(token.TypeLT, "<"),
				token.New(token.TypeInt, "10"),
				token.New(token.TypeRParen, ")"),
				token.New(token.TypeLBrace, "{"),
				token.New(token.TypeReturn, "return"),
				token.New(token.TypeBool, "true"),
				token.New(token.TypeSemicolon, ";"),
				token.New(token.TypeRBrace, "}"),
				token.New(token.TypeElse, "else"),
				token.New(token.TypeLBrace, "{"),
				token.New(token.TypeReturn, "return"),
				token.New(token.TypeBool, "false"),
				token.New(token.TypeSemicolon, ";"),
				token.New(token.TypeRBrace, "}"),
			},
		},
		"Equality operators": {
			input: `
10 == 10;
10 != 9;`,
			toks: []token.Token{
				token.New(token.TypeInt, "10"),
				token.New(token.TypeEQ, "=="),
				token.New(token.TypeInt, "10"),
				token.New(token.TypeSemicolon, ";"),
				token.New(token.TypeInt, "10"),
				token.New(token.TypeNotEQ, "!="),
				token.New(token.TypeInt, "9"),
				token.New(token.TypeSemicolon, ";"),
			},
		},
		"Strings": {
			input: `
"foobar"
"foo bar"
`,
			toks: []token.Token{
				token.New(token.TypeString, "foobar"),
				token.New(token.TypeString, "foo bar"),
			},
		},
		"Slices": {
			input: "[]; [1]; [1,2];",
			toks: []token.Token{
				token.New(token.TypeLBrakt, "["),
				token.New(token.TypeRBrakt, "]"),
				token.New(token.TypeSemicolon, ";"),
				token.New(token.TypeLBrakt, "["),
				token.New(token.TypeInt, "1"),
				token.New(token.TypeRBrakt, "]"),
				token.New(token.TypeSemicolon, ";"),
				token.New(token.TypeLBrakt, "["),
				token.New(token.TypeInt, "1"),
				token.New(token.TypeComma, ","),
				token.New(token.TypeInt, "2"),
				token.New(token.TypeRBrakt, "]"),
				token.New(token.TypeSemicolon, ";"),
				token.New(token.TypeEOF, ""),
			},
		},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New(tc.input)
			for _, want := range tc.toks {
				got := l.NextToken()
				is.Equal(t, want.Type(), got.Type())
				is.Equal(t, want.Literal(), got.Literal())
			}
		})
	}
}
