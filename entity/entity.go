// package entity is like the ast package to the parser package. It defines the
// entities or objects that the mmm programming language can evaluate. Therefore
// if you want to evaluate a function there's an entity.Fn for that, and if you
// want to evaluate some other entity or object or type or whatever you want to
// call it, there will need to be a representation in the ast package first.
//
// If you look hard enough the `Inspect` method on the E interface looks an
// awful lot like a `String` method for a [fmt.Stringer] and you'd be correct.
// You'd also see that all ast.XXXX structs have a `String` method, so
// theoretically, we should be able to take the AST representation and use it as
// the `Inspect` method for all [E]s. A small problem is not all
// ast.XXXX.String() match up directly with what we'd like to see when printing
// to the user. And a big problem is we have no access to those [ast.Node]s in
// this package.
package entity

import (
	"bytes"
	"strconv"
	"strings"

	"mmm/ast"
)

// E is an Entity that satisfies having a type in the mmm language and has a
// representation.
type E interface {
	Type() Type
	Inspect() string
}

type BuiltinFn func(...E) E

type Env struct {
	parent *Env
	store map[string]E
}

func NewEnv() Env {
	return Env{store: map[string]E{}}
}

func NewEnvWith(parent *Env) Env {
	return Env{store: map[string]E{}, parent: parent}
}

func (e Env) Get(name string) (E, bool) {
	v, ok := e.store[name]
	if !ok && e.parent != nil {
		v, ok = e.parent.Get(name)
	}
	return v, ok
}

func (e Env) Set(name string, val E) E {
	e.store[name] = val
	return val
}

type Type uint8

const (
	TypeError Type = iota
	TypeNull
	TypeInt
	TypeBool
	TypeReturn
	TypeFn
	TypeString
	TypeBuiltin
	TypeSlice
)

func (t Type) String() string {
	switch t {
	case TypeError:
		return "Error"
	case TypeNull:
		return "Null"
	case TypeInt:
		return "Int"
	case TypeBool:
		return "Bool"
	case TypeReturn:
		return "Return"
	case TypeFn:
		return "Fn"
	case TypeString:
		return "String"
	case TypeBuiltin:
		return "Builtin"
	case TypeSlice:
		return "Slice"
	default:
		return "Unknown"
	}
}


type Int struct {
	Value int64
}

func (Int) Type() Type { return TypeInt }
func (i Int) Inspect() string { return strconv.Itoa(int(i.Value)) }

type Bool struct {
	Value bool
}

func (Bool) Type() Type { return TypeBool }
func (b Bool) Inspect() string { return strconv.FormatBool(b.Value) }

type Null struct{}

func (Null) Type() Type { return TypeNull }
func (Null) Inspect() string { return "null" }

type Return struct {
	Value E
}

func (Return) Type() Type { return TypeReturn }
func (r Return) Inspect() string { return r.Value.Inspect() }

type Error struct {
	Message string
}

func (Error) Type() Type { return TypeError }
func (e Error) Inspect() string { return "ERROR: " + e.Message }

type Fn struct {
	Params []ast.Ident
	Body ast.BlockStmt
	Env Env
}

func (Fn) Type() Type { return TypeFn }
func (f Fn) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Params {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())
	return out.String()
}

type String struct {
	Value string
}

func (String) Type() Type { return TypeString }
func (s String) Inspect() string { return s.Value }

type Builtin struct {
	Fn BuiltinFn
}

func (Builtin) Type() Type { return TypeBuiltin }
func (Builtin) Inspect() string { return "builtin function" }

type Slice struct {
	Values []E
}

func (Slice) Type() Type { return TypeSlice }
func (s Slice) Inspect() string {
	var out bytes.Buffer
	values := make([]string, len(s.Values))
	for i, v := range s.Values {
		values[i] = v.Inspect()
	}
	out.WriteString("[")
	out.WriteString(strings.Join(values, ", "))
	out.WriteString("]")
	return out.String()
}
