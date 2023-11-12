package parser

import (
	"errors"
	"mmm/ast"
	"mmm/lexer"
	"mmm/token"
	"strconv"
)

type priority uint8

const (
	_ priority = iota
	priorityLowest
	priorityEquals // priorityEquals e.g. x == x
	priorityLessGreater // priorityLessGreater e.g. < or >
	prioritySum
	priorityProduct
	priorityPrefix // priorityPrefix e.g. -x, !x
	priorityCall // priorityCall e.g. call(x)
	priorityIndex // priorityIndex e.g. slice[1]
)

// These are part of Pratt parsing implementations to give priority (or
// precedence) to certain token types. With more priority expressions are
// evaluated first. This is required for including things like PEMDAS in
// arithmatic.
type (
	prefixParseFunc func() ast.Expr
	infixParseFunc func(ast.Expr) ast.Expr
)

type TokenError struct {
	want token.Type
	got  token.Type
}

func (e TokenError) Error() string {
	return "expected next token to be " + e.want.String() +
		", got " + e.got.String()
}

// Parser parses tokens passed to it by the [*lexer.Lexer] and ultimately
// returns a valid mmm program.
type Parser struct {
	l *lexer.Lexer
	// ctok is the current token the Parser is looking at.
	ctok token.Token
	// ntok is the next token the Parser will consume. It's needed to have this in
	// case ctok doesn't give us enough information on how to build the AST.
	ntok token.Token
	// errs is all of the errors found while trying to parse statements to a
	// program.
	errs []error
	prefixes func(token.Type)prefixParseFunc
	infixes func(token.Type)infixParseFunc
	priorities func(token.Type)priority
}

// New returns a [*Parser] that is ready to make an [*ast.Program] with the
// passed in [*lexer.Lexer].
func New(l *lexer.Lexer) *Parser {
	p := &Parser{ l: l }
	p.prefixes = func(t token.Type) prefixParseFunc {
		switch t {
		case token.TypeIdent:
			return func() ast.Expr {

				return ast.NewIdent(p.ctok.Literal())
			}
		case token.TypeInt:
			return func() ast.Expr {
				v, err := strconv.ParseInt(p.ctok.Literal(), 0,64)
				if err != nil {
					p.errs = append(p.errs, errors.New("could not parse "+
					p.ctok.Type().String() + " as integer"))
					return nil
				}
				return ast.NewInteger(v)
			}
		case token.TypeBang, token.TypeMinus:
			return func() ast.Expr {
				t := p.ctok
				p.nextToken()
				return ast.NewPrefixExpr(t, t.Literal(), p.parseExpression(priorityPrefix))
			}
		case token.TypeBool:
			return func() ast.Expr {
				return ast.NewBool(p.ctok.Literal() == "true")
			}
		case token.TypeLParen:
			return func() ast.Expr {
				p.nextToken()
				expr := p.parseExpression(priorityLowest)
				if !p.peek(token.TypeRParen) {
					return nil
				}
				return expr
			}
		case token.TypeIf:
			return func() ast.Expr { // if 
				if !p.peek(token.TypeLParen) { // (
					return nil
				}
				p.nextToken()
				// x == y
				cond := p.parseExpression(priorityLowest)
				if !p.peek(token.TypeRParen) { // )
					return nil
				}
				if !p.peek(token.TypeLBrace) { // {
					return nil
				}
				consq := p.parseBlock() // ... }
				var alt *ast.BlockStmt
				if p.ntok.Type() == token.TypeElse {
					p.nextToken()
					if !p.peek(token.TypeLBrace) { // {
						return nil
					}
					blk := p.parseBlock()
					alt = &blk
				}
				return ast.NewIfExpr(cond, consq, alt)
			}
		case token.TypeFn:
			return func() ast.Expr {
				if !p.peek(token.TypeLParen) {
					return nil
				}
				params := p.parseFnParams()
				if !p.peek(token.TypeLBrace) {
					return nil
				}
				return ast.NewFunction(params, p.parseBlock())
			}
		case token.TypeString:
			return func() ast.Expr {
				return ast.NewString(p.ctok.Literal())
			}
		case token.TypeLBrakt:
			return func() ast.Expr {
				return ast.NewSlice(p.parseExprSlice(token.TypeRBrakt)...)
			}
		default:
			return nil
		}
	}
	p.infixes = func(t token.Type) infixParseFunc {
		switch t {
		case token.TypeEQ, token.TypeNotEQ, token.TypeLT, token.TypeGT,
		token.TypePlus, token.TypeMinus, token.TypeSlash, token.TypeStar:
			return func(left ast.Expr) ast.Expr {
				t, pri := p.ctok, p.priorities(p.ctok.Type())
				p.nextToken()
				return ast.NewInfixExpr(t, t.Literal(), left, p.parseExpression(pri))
			}
		case token.TypeLParen:
			// This isn't actually an infix operator, but if fits nicely with what we
			// would like to see from a call, e.g. blah(1, 2, 3, 4)
			return func(fn ast.Expr) ast.Expr {
				return ast.NewCallExpr(fn, p.parseExprSlice(token.TypeRParen))
			}
		case token.TypeLBrakt:
			// This isn't actually an infix operator, but if fits nicely with what we
			// would like to see from an index into a slice, e.g. blah[1]
			return func(e ast.Expr) ast.Expr {
				p.nextToken()
				idx := p.parseExpression(priorityLowest)
				if !p.peek(token.TypeRBrakt) {
					return nil
				}
				return ast.NewIndex(e, idx)
			}
		default:
			return nil
		}
	}
	p.priorities = func(t token.Type) priority {
		switch t {
		case token.TypeEQ, token.TypeNotEQ:
			return priorityEquals
		case token.TypeLT, token.TypeGT:
			return priorityLessGreater
		case token.TypePlus, token.TypeMinus:
			return prioritySum
		case token.TypeSlash, token.TypeStar:
			return priorityProduct
		case token.TypeLParen:
			return priorityCall
		case token.TypeLBrakt:
			return priorityIndex
		default:
			return priorityLowest
		}
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p Parser) Errors() []string {
	s := make([]string, len(p.errs))
	for i, e := range p.errs {
		s[i] = e.Error()
	}
	return s
}

// Parse creates a [*ast.Program] out of the tokens the Parser has.
func (p *Parser) Parse() ast.Program {
	var program ast.Program
	for p.ctok.Type() != token.TypeEOF {
		if s := p.parseStatement(); s != nil {
			program.Statements = append(program.Statements, s)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.ctok.Type() {
	case token.TypeLet:
		if !p.peek(token.TypeIdent) {
			return nil
		}
		id := ast.NewIdent(p.ctok.Literal())
		if !p.peek(token.TypeAssign) {
			return nil
		}
		p.nextToken()
		expr := p.parseExpression(priorityLowest)
		if p.ntok.Type() == token.TypeSemicolon {
			p.nextToken()
		}
		return ast.NewLetStmt(id, expr)
	case token.TypeReturn:
		p.nextToken()
		expr := p.parseExpression(priorityLowest)
		if p.ntok.Type() == token.TypeSemicolon {
			p.nextToken()
		}
		return ast.NewRetStmt(expr)
	default:
		es := ast.NewExprStmt(p.ctok, p.parseExpression(priorityLowest))
		if p.ntok.Type() == token.TypeSemicolon {
			p.nextToken()
		}
		return es
	}
}

func (p *Parser) parseExpression(pr priority) ast.Expr {
	prefix := p.prefixes(p.ctok.Type())
	if prefix == nil {
		p.errs = append(p.errs, errors.New("no prefix parse function for " +
p.ctok.Type().String() + " found"))
		return nil
	}
	left := prefix()
	for p.ntok.Type() != token.TypeSemicolon && pr < p.priorities(p.ntok.Type()) {
		infix := p.infixes(p.ntok.Type())
		if infix == nil {
			return left
		}
		p.nextToken()
		left = infix(left)
	}
	return left
}

// parseBlock consumes a block statment defined as { ... }.
func (p *Parser) parseBlock() ast.BlockStmt {
	p.nextToken()
	var ss []ast.Statement
	for p.ctok.Type() != token.TypeRBrace && p.ctok.Type() != token.TypeEOF {
		if s := p.parseStatement(); s != nil {
			ss = append(ss, s)
		}
		p.nextToken()
	}
	return ast.NewBlockStmt(ss)
}

func (p *Parser) parseFnParams() []ast.Ident {
	var idents []ast.Ident
	if p.ntok.Type() == token.TypeRParen {
		p.nextToken()
		return idents
	}
	p.nextToken() // x
	// (x, y)
	idents = append(idents, ast.NewIdent(p.ctok.Literal()))
	for p.ntok.Type() == token.TypeComma {
		p.nextToken() // ,
		p.nextToken() // y
		idents = append(idents, ast.NewIdent(p.ctok.Literal()))
	}
	if !p.peek(token.TypeRParen) {
		return nil
	}
	return idents
}

func (p *Parser) nextToken() {
	p.ctok = p.ntok
	p.ntok = p.l.NextToken()
}

func (p *Parser) peek(t token.Type) bool {
	if got := p.ntok.Type(); got != t {
		p.errs = append(p.errs, TokenError{want: t, got: got})
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) parseExprSlice(end token.Type) []ast.Expr {
	if p.ntok.Type() == token.TypeRParen {
		p.nextToken()
		return nil
	}
	// (x, y)
	p.nextToken() // (
	var vals []ast.Expr
	vals = append(vals, p.parseExpression(priorityLowest)) // x
	for p.ntok.Type() == token.TypeComma {
		p.nextToken() // ,
		p.nextToken() // y
		vals = append(vals, p.parseExpression(priorityLowest))
	}
	if !p.peek(end) {
		return nil
	}
	return vals
}
