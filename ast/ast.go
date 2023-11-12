package ast

import (
	"bytes"
	"fmt"
	"mmm/token"
	"strconv"
	"strings"
)

type Node interface {
	TokenLiteral() string
	fmt.Stringer
}

type Statement interface {
	Node
	isStmt()
}

type Expr interface {
	Node
	isExpr()
}

type Program struct {
	Statements []Statement
}

func (p Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}
	return p.Statements[0].TokenLiteral()
}

func (p Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStmt struct {
	t     token.Token
	name  Ident
	value Expr
}

func NewLetStmt(id Ident, value Expr) LetStmt {
	return LetStmt{
		t:     token.New(token.TypeLookup, "let"),
		name:  id,
		value: value,
	}
}

func (LetStmt) isStmt()                {}
func (l LetStmt) TokenLiteral() string { return l.t.Literal() }
func (l LetStmt) Name() string         { return l.name.value }
func (l LetStmt) Value() Expr        { return l.value }
func (ls LetStmt) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.name.String())
	out.WriteString(" = ")
	if ls.value != nil {
		out.WriteString(ls.value.String())
	}
	out.WriteString(";")
	return out.String()
}

// RetStmt is the Node representing a return statment.
type RetStmt struct {
	t     token.Token
	value Expr
}

func NewRetStmt(value Expr) RetStmt {
	return RetStmt{
		t:     token.New(token.TypeLookup, "return"),
		value: value,
	}
}

func (RetStmt) isStmt()                {}
func (rs RetStmt) TokenLiteral() string { return rs.t.Literal() }
func (rs RetStmt) Value() Expr { return rs.value }
func (rs RetStmt) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.value != nil {
		out.WriteString(rs.value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExprStmt is required for mmm because it is legal to do this:
//
//	let foo = 1;
//	foo + 20;
//
// which is just an expression as a statement if the user wants to print the
// value of just the expression.
type ExprStmt struct {
	t     token.Token
	value Expr
}

func NewExprStmt(t token.Token, v Expr) ExprStmt {
	return ExprStmt{t: t, value: v}
}

func (ExprStmt) isStmt()                 {}
func (es ExprStmt) TokenLiteral() string { return es.t.Literal() }
func (es ExprStmt) Expression() Expr {
	return es.value
}
func (es ExprStmt) String() string {
	if es.value == nil {
		return ""
	}
	return es.value.String()
}

// Ident is the literal identity name stored in a value
// e.g. the `foobar` in let foobar = 5;
type Ident struct {
	t     token.Token
	value string
}

func NewIdent(value string) Ident {
	return Ident{t: token.New(token.TypeIdent, value), value: value}
}

func (Ident) isExpr()                {}
func (i Ident) TokenLiteral() string { return i.t.Literal() }
func (i Ident) String() string       { return i.value }

type Integer struct {
	t     token.Token
	value int64
}

func NewInteger(v int64) Integer {
	return Integer{t: token.New(token.TypeInt, strconv.Itoa(int(v))), value: v}
}

func (Integer) isExpr()                {}
func (i Integer) TokenLiteral() string { return i.t.Literal() }
func (i Integer) Value() int64         { return i.value }
func (i Integer) String() string       { return i.t.Literal() }

type PrefixExpr struct {
	t     token.Token
	op    string
	right Expr
}

func NewPrefixExpr(t token.Token, operator string, right Expr) PrefixExpr {
	return PrefixExpr{
		t:     t,
		op:    operator,
		right: right,
	}
}

func (PrefixExpr) isExpr()                 {}
func (pe PrefixExpr) TokenLiteral() string { return pe.t.Literal() }
func (pe PrefixExpr) Operator() string     { return pe.op }
func (pe PrefixExpr) Right() Expr          { return pe.right }
func (pe PrefixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.op)
	out.WriteString(pe.right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpr struct {
	t     token.Token
	left  Expr
	op    string
	right Expr
}

func NewInfixExpr(t token.Token, operator string, left, right Expr) InfixExpr {
	return InfixExpr{
		t:     t,
		left:  left,
		op:    operator,
		right: right,
	}
}

func (InfixExpr) isExpr()                 {}
func (ie InfixExpr) TokenLiteral() string { return ie.t.Literal() }
func (ie InfixExpr) Left() Expr           { return ie.left }
func (ie InfixExpr) Right() Expr          { return ie.right }
func (ie InfixExpr) Operator() string     { return ie.op }
func (ie InfixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.left.String())
	out.WriteString(" " + ie.op + " ")
	out.WriteString(ie.right.String())
	out.WriteString(")")
	return out.String()
}

type Bool struct {
	t     token.Token
	value bool
}

func NewBool(value bool) Bool {
	return Bool{t: token.New(token.TypeLookup, strconv.FormatBool(value)), value: value}
}

func (Bool) isExpr()                {}
func (b Bool) TokenLiteral() string { return b.t.Literal() }
func (b Bool) String() string       { return b.t.Literal() }
func (b Bool) Value() bool          { return b.value }

type BlockStmt struct {
	t          token.Token
	Statements []Statement
	valid      bool
}

func NewBlockStmt(s []Statement) BlockStmt {
	return BlockStmt{
		t:          token.New(token.TypeLBrace, "{"),
		Statements: s,
		valid:      true,
	}
}

func (BlockStmt) isStmt()                 {}
func (bs BlockStmt) TokenLiteral() string { return bs.t.Literal() }
func (bs BlockStmt) OK() bool             { return bs.valid }
func (bs BlockStmt) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type IfExpr struct {
	t           token.Token
	Condition   Expr
	Consequence BlockStmt
	Alternative BlockStmt
}

func NewIfExpr(
	condition Expr, consequence BlockStmt, alternative *BlockStmt,
) IfExpr {
	var alt BlockStmt
	if alternative != nil {
		alt = *alternative
	}
	return IfExpr{
		t:           token.New(token.TypeLookup, "if"),
		Condition:   condition,
		Consequence: consequence,
		Alternative: alt,
	}
}

func (IfExpr) isExpr()                 {}
func (ie IfExpr) TokenLiteral() string { return ie.t.Literal() }
func (ie IfExpr) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative.valid {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type Function struct {
	t      token.Token
	Params []Ident
	Body   BlockStmt
}

func NewFunction(params []Ident, body BlockStmt) Function {
	return Function{
		t:      token.New(token.TypeLookup, "fn"),
		Params: params,
		Body:   body,
	}
}

func (Function) isExpr()                {}
func (f Function) TokenLiteral() string { return f.t.Literal() }
func (f Function) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Params {
		params = append(params, p.String())
	}
	out.WriteString(f.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())
	return out.String()
}

type CallExpr struct {
	t    token.Token // t is '(' token
	Fn   Expr
	Args []Expr
}

func NewCallExpr(fn Expr, args []Expr) CallExpr {
	return CallExpr{
		t:    token.New(token.TypeLParen, "("),
		Fn:   fn,
		Args: args,
	}
}

func (CallExpr) isExpr()                 {}
func (ce CallExpr) TokenLiteral() string { return ce.t.Literal() }
func (ce CallExpr) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Args {
		args = append(args, a.String())
	}
	out.WriteString(ce.Fn.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

type String struct {
	t token.Token
	value string
}

func NewString(value string) String {
	return String{t: token.New(token.TypeString, value), value: value}
}

func (String) isExpr() {}
func (s String) TokenLiteral() string { return s.t.Literal() }
func (s String) String() string { return s.value }

type Slice struct {
	t token.Token
	values []Expr
}

func NewSlice(values ...Expr) Slice {
	return Slice{t: token.New(token.TypeLBrakt, "["), values: values}
}

func (Slice) isExpr() {}
func (a Slice) TokenLiteral() string { return a.t.Literal() }
func (a Slice) Values() []Expr { return a.values }
func (a Slice) String() string {
	var out bytes.Buffer
	values := make([]string, len(a.values))
	for i, v := range a.values {
		values[i] = v.String()
	}
	out.WriteString("[")
	out.WriteString(strings.Join(values, ", "))
	out.WriteString("]")
	return out.String()
}

type Index struct {
	t token.Token
	left Expr
	idx Expr
}

func NewIndex(left, idx Expr) Index {
	return Index{t: token.New(token.TypeLBrakt, "["), left: left, idx: idx}
}

func (Index) isExpr() {}
func (i Index) TokenLiteral() string { return i.t.Literal() }
func (i Index)	Left() Expr { return i.left }
func (i Index)	Idx() Expr { return i.idx }
func (i Index) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.left.String())
	out.WriteString("[")
	out.WriteString(i.idx.String())
	out.WriteString("])")
	return out.String()
}

