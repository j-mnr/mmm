package eval

import (
	"fmt"
	"mmm/ast"
	"mmm/entity"
)

var (
	builtins = map[string]entity.Builtin{
		"len": {Fn: func(e ...entity.E) entity.E {
			if len(e) != 1 {
				return newErr("len only accepts one argument.")
			}
			switch v := e[0].(type) {
			case entity.String:
				return entity.Int{Value: int64(len(v.Value))}
			case entity.Slice:
				return entity.Int{Value: int64(len(v.Values))}
			default:
				return newErr("argument to `len` not supported, got %s", v.Type())
			}
			},
		},
	}
	_true = entity.Bool{Value: true}
	_false = entity.Bool{Value: false}
	null = entity.Null{}
)

func Eval(node ast.Node, env entity.Env) entity.E {
	switch node := node.(type) {
	// Statements
	case ast.Program:
		return evalProgram(node.Statements, env)
	case ast.ExprStmt:
		return Eval(node.Expression(), env)
	case ast.BlockStmt:
		return evalBlock(node, env)
	case ast.RetStmt:
		v := Eval(node.Value(), env)
		if isErr(v) {
			return v
		}
		return entity.Return{Value: v}
	case ast.LetStmt:
		val := Eval(node.Value(), env)
		if isErr(val) {
			return val
		}
		env.Set(node.Name(), val)
		return nil
	// Expressions
	case ast.Bool:
		if node.Value() {
			return _true
		}
		return _false
	case ast.Integer:
		return entity.Int{Value: node.Value()}
	case ast.PrefixExpr:
		v := Eval(node.Right(), env)
		if isErr(v) {
			return v
		}
		return evalPrefix(node.Operator(), v)
	case ast.InfixExpr:
		l := Eval(node.Left(), env)
		if isErr(l) {
			return l
		}
		r := Eval(node.Right(), env)
		if isErr(r) {
			return r
		}
		return evalInfix(l, node.Operator(), r)
	case ast.IfExpr:
		return evalIf(node, env)
	case ast.Ident:
		if val, ok := env.Get(node.String()); ok {
			return val
		}
		if b, ok := builtins[node.String()]; ok {
			return b
		}
		return newErr("identifier not found: " + node.String())
	case ast.Function:
		return entity.Fn{Env: env, Params: node.Params, Body: node.Body}
	case ast.CallExpr:
		fn := Eval(node.Fn, env)
		if isErr(fn) {
			return fn
		}
		args := evalExpressions(node.Args, env)
		if len(args) == 1 && isErr(args[0]) {
			return args[0]
		}
		return evalFn(fn, args)
	case ast.String:
		return entity.String{Value: node.String()}
	case ast.Slice:
		vals := evalExpressions(node.Values(), env)
		if len(vals) == 1 && isErr(vals[0]) {
			return vals[0]
		}
		return entity.Slice{Values: vals}
	case ast.Index:
		l := Eval(node.Left(), env)
		if isErr(l) {
			return l
		}
		i := Eval(node.Idx(), env)
		if isErr(i) {
			return i
		}
		return evalIndex(l, i)
	default:
		return nil
	}
}

func staticBool(isTrue bool) entity.Bool {
	if isTrue {
		return _true
	}
	return _false
}

func evalProgram(ss []ast.Statement, env entity.Env) entity.E {
	var res entity.E
	for _, s := range ss {
		res = Eval(s, env)
		switch res := res.(type) {
		case entity.Return:
			return res.Value
		case entity.Error:
			return res
		}
	}
	return res
}

func evalPrefix(op string, right entity.E) entity.E {
	switch op {
	case "!":
		switch right {
		case _true:
			return _false
		case _false:
			return _true
		case null:
			return _true
		default:
			return _false
		}
	case "-":
		if right.Type() != entity.TypeInt {
			return newErr("unknown operator: -%s", right.Type())
		}
		return entity.Int{Value: -right.(entity.Int).Value}
	default:
		return newErr("unknown operator: %s%s", op, right.Type())
	}
}

func evalInfix(left entity.E, op string, right entity.E) entity.E {
	if left.Type() != right.Type() {
		return newErr("type mismatch: %s %s %s", left.Type(), op, right.Type())
	}
	if left.Type() == entity.TypeInt && right.Type() == entity.TypeInt {
		return evalIntInfix(left, op, right)
	}
	if left.Type() == entity.TypeString && right.Type() == entity.TypeString {
	if op != "+" {
		return newErr("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
	return entity.String{Value: left.(entity.String).Value + right.(entity.String).Value}
	}
	switch op {
	case "==":
		return staticBool(left == right)
	case "!=":
		return staticBool(left != right)
	default:
		return newErr("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalIntInfix(left entity.E, op string, right entity.E) entity.E {
	lval, rval := left.(entity.Int).Value, right.(entity.Int).Value
	switch op {
	case "+":
		return entity.Int{Value: lval + rval}
	case "-":
		return entity.Int{Value: lval - rval}
	case "/":
		return entity.Int{Value: lval / rval}
	case "*":
		return entity.Int{Value: lval * rval}
	case "<":
		return staticBool(lval < rval)
	case ">":
		return staticBool(lval > rval)
	case "==":
		return staticBool(lval == rval)
	case "!=":
		return staticBool(lval != rval)
	default:
		return newErr("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalIf(ife ast.IfExpr, env entity.Env) entity.E {
	c := Eval(ife.Condition, env)
	if isErr(c) {
		return c
	}
	switch {
	case isTruthy(c):
		return Eval(ife.Consequence, env)
	case ife.Alternative.OK():
		return Eval(ife.Alternative, env)
	default:
		return null
	}
}

func isTruthy(e entity.E) bool {
	switch e {
	case null:
		return false
	case _true:
		return true
	case _false:
		return false
	default:
		return true
	}
}

func evalBlock(blk ast.BlockStmt, env entity.Env) entity.E {
	var res entity.E
	for _, s := range blk.Statements {
		res = Eval(s, env)
		if res != nil  {
			switch res.Type() {
			case entity.TypeReturn, entity.TypeError:
				return res
			}
		}
	}
	return res
}

func newErr(format string, a ...any) entity.Error {
	return entity.Error{Message: fmt.Sprintf(format, a...)}
}

func isErr(e entity.E) bool {
	if e == nil {
		return false
	}
	return e.Type() == entity.TypeError
}

func evalExpressions(exprs []ast.Expr, env entity.Env) []entity.E {
	var res []entity.E
	for _, e := range exprs {
		val := Eval(e, env)
		if isErr(val) {
			return []entity.E{val}
		}
		res = append(res, val)
	}
	return res
}

func evalFn(fn entity.E, args []entity.E) entity.E {
	switch fn := fn.(type) {
	case entity.Fn:
	env := entity.NewEnvWith(&fn.Env)
	for i, p := range fn.Params {
		env.Set(p.String(), args[i])
	}
	val := Eval(fn.Body, env)
	if ret, ok := val.(entity.Return); ok {
		return ret.Value
	}
	return val
	case entity.Builtin:
		return fn.Fn(args...)
	default:
		return newErr("not a function: %s", fn.Type())
	}
}

func evalIndex(left, idx entity.E) entity.E {
	switch left.Type() {
	case entity.TypeSlice:
			left := left.(entity.Slice)
			if idx.Type() != entity.TypeInt {
				return null
			}
			i := idx.(entity.Int).Value
			if i < 0 || i > int64(len(left.Values)-1) {
				return null
			}
			return left.Values[i]
		default:
			return newErr("index operator not supported for %s", left.Type())
	}
}
