package parser_test

import (
	"mmm/ast"
	"mmm/is"
	"mmm/lexer"
	"mmm/parser"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	t.Parallel()
	t.Run("Let Statements", func(t *testing.T) {
		t.Parallel()
		p := parser.New(lexer.New(`
let x = 5;
let y = 10;
let foobar = 9320812;`))
		program := p.Parse()
		checkErrors(t, p.Errors())
		for i, want := range []string{"x", "y", "foobar"} {
			stmt := program.Statements[i]
			is.Equal(t, "let", stmt.TokenLiteral())
			let := stmt.(ast.LetStmt)
			is.Equal(t, want, let.Name())
			// is.Equal(t, want, let.Name.TokenLiteral())
		}
	})
	t.Run("Return Statements", func(t *testing.T) {
		t.Parallel()
		p := parser.New(lexer.New(`
return 5;
return 10;
return 9320812;`))
		program := p.Parse()
		if len(program.Statements) == 0 {
			t.Fatal("The program was nil")
		}
		checkErrors(t, p.Errors())
	})
	t.Run("Identifier Expression", func(t *testing.T) {
		t.Parallel()
		p := parser.New(lexer.New(`foobar;`))
		program := p.Parse()
		checkErrors(t, p.Errors())
		stmt := program.Statements[0].(ast.ExprStmt)
		ident := stmt.Expression().(ast.Ident)
		is.Equal(t, "foobar", ident.TokenLiteral())
	})
	t.Run("Integer Literal", func(t *testing.T) {
		t.Parallel()
		p := parser.New(lexer.New(`5;`))
		program := p.Parse()
		checkErrors(t, p.Errors())
		stmt := program.Statements[0].(ast.ExprStmt)
		ident := stmt.Expression().(ast.Integer)
		is.Equal(t, "5", ident.TokenLiteral())
	})
	t.Run("Parse Prefix Int Expressions", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input    string
			operator string
			want     int64
		}{
			"Bang":     {input: "!5;", operator: "!", want: 5},
			"Negative": {input: "-15;", operator: "-", want: 15},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				p := parser.New(lexer.New(tc.input))
				program := p.Parse()
				checkErrors(t, p.Errors())
				is.Equal(t, 1, len(program.Statements))
				stmt := program.Statements[0].(ast.ExprStmt)
				prefix := stmt.Expression().(ast.PrefixExpr)
				is.Equal(t, tc.operator, prefix.Operator())
				is.Equal(t, tc.want, prefix.Right().(ast.Integer).Value())
			})
		}
	})
	t.Run("Parse Prefix Bool Expressions", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input    string
			operator string
			want     bool
		}{
			"True":  {input: "!true;", operator: "!", want: true},
			"False": {input: "!false;", operator: "!", want: false},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				p := parser.New(lexer.New(tc.input))
				program := p.Parse()
				checkErrors(t, p.Errors())
				is.Equal(t, 1, len(program.Statements))
				stmt := program.Statements[0].(ast.ExprStmt)
				prefix := stmt.Expression().(ast.PrefixExpr)
				is.Equal(t, tc.operator, prefix.Operator())
				is.Equal(t, tc.want, prefix.Right().(ast.Bool).Value())
			})
		}
	})
	t.Run("Parse Infix Integer Expressions", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input    string
			left     int64
			operator string
			right    int64
		}{
			"Plus":         {input: "5 + 5;", left: 5, operator: "+", right: 5},
			"Minus":        {input: "5 - 5;", left: 5, operator: "-", right: 5},
			"Star":         {input: "5 * 5;", left: 5, operator: "*", right: 5},
			"Slash":        {input: "5 / 5;", left: 5, operator: "/", right: 5},
			"Greater Than": {input: "5 > 5;", left: 5, operator: ">", right: 5},
			"Less Than":    {input: "5 < 5;", left: 5, operator: "<", right: 5},
			"Equal":        {input: "5 == 5;", left: 5, operator: "==", right: 5},
			"Not Equal":    {input: "5 != 5;", left: 5, operator: "!=", right: 5},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				p := parser.New(lexer.New(tc.input))
				program := p.Parse()
				checkErrors(t, p.Errors())
				is.Equal(t, 1, len(program.Statements))
				stmt := program.Statements[0].(ast.ExprStmt)
				infix := stmt.Expression().(ast.InfixExpr)
				is.Equal(t, tc.left, infix.Left().(ast.Integer).Value())
				is.Equal(t, tc.operator, infix.Operator())
				is.Equal(t, tc.right, infix.Right().(ast.Integer).Value())
			})
		}
	})
	t.Run("Parse Infix Bool Expressions", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input    string
			left     bool
			operator string
			right    bool
		}{
			"True is True":      {input: "true == true", left: true, operator: "==", right: true},
			"True is not False": {input: "true != false", left: true, operator: "!=", right: false},
			"False is False":    {input: "false == false;", left: false, operator: "==", right: false},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				p := parser.New(lexer.New(tc.input))
				program := p.Parse()
				if len(program.Statements) == 0 {
					t.Fatal("The program was nil")
				}
				checkErrors(t, p.Errors())
				is.Equal(t, 1, len(program.Statements))
				stmt := program.Statements[0].(ast.ExprStmt)
				infix := stmt.Expression().(ast.InfixExpr)
				is.Equal(t, tc.left, infix.Left().(ast.Bool).Value())
				is.Equal(t, tc.operator, infix.Operator())
				is.Equal(t, tc.right, infix.Right().(ast.Bool).Value())
			})
		}
	})
	t.Run("Operator Priority Parsing", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input, want string
		}{
			"Prefix with Product": {
				input: "-a * b", want: "((-a) * b)",
			},
			"Two prefixes": {
				input: "!-a",
				want:  "(!(-a))",
			},
			"Multiple Sum": {
				input: "a + b + c",
				want:  "((a + b) + c)",
			},
			"Multiple Product": {
				input: "a * b * c",
				want:  "((a * b) * c)",
			},
			"Multiple Priorities": {
				input: "a + b * c + d / e - f",
				want:  "(((a + (b * c)) + (d / e)) - f)",
			},
			"Multiple Statements": {
				input: "3 + 4; -5 * 5",
				want:  "(3 + 4)((-5) * 5)",
			},
			"Boolean operators": {
				input: "5 > 4 == 3 < 4 != 8",
				want:  "(((5 > 4) == (3 < 4)) != 8)",
			},
			"All priorities": {
				input: "5 * 3 + 2 > 1 == 1 < 2 * 3 + 5",
				want:  "((((5 * 3) + 2) > 1) == (1 < ((2 * 3) + 5)))",
			},
			"Groups take priority": {
				input: "-((5 + 5) * 5)",
				want:  "(-((5 + 5) * 5))",
			},
			"Calls are highest priority": {
				input: "add(a, b, 1, 2, 3 * 4, add(1 / 1, 2))",
				want:  "add(a, b, 1, 2, (3 * 4), add((1 / 1), 2))",
			},
			"Slices": {
				input: "add(a, b[1], 1, 2, b[3 * 4], add(1 / [1, 2][1]))",
				want:  "add(a, (b[1]), 1, 2, (b[(3 * 4)]), add((1 / ([1, 2][1]))))",
			},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				p := parser.New(lexer.New(tc.input))
				program := p.Parse()
				checkErrors(t, p.Errors())
				is.Equal(t, tc.want, program.String())
			})
		}
	})
	t.Run("If/Else Expression", func(t *testing.T) {
		t.Parallel()
		p := parser.New(lexer.New("if (x < y) { x } else { y }"))
		program := p.Parse()
		checkErrors(t, p.Errors())
		is.Equal(t, 1, len(program.Statements))
		expr := program.Statements[0].(ast.ExprStmt).Expression().(ast.IfExpr)

		infix := expr.Condition.(ast.InfixExpr)
		is.Equal(t, "x", infix.Left().String())
		is.Equal(t, "<", infix.Operator())
		is.Equal(t, "y", infix.Right().String())

		is.Equal(t, 1, len(expr.Consequence.Statements))
		stmt := expr.Consequence.Statements[0].(ast.ExprStmt)
		is.Equal(t, "x", stmt.Expression().String())
		is.Equal(t, true, expr.Alternative.OK())
		is.Equal(t, "y", expr.Alternative.Statements[0].String())
	})
	t.Run("Functions", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input      string
			wantParams []string
			wantBody   string
		}{
			"Zero": {
				input:      "fn () {};",
				wantParams: nil,
				wantBody:   "",
			},
			"One": {
				input:      "fn(x) { return x; }",
				wantParams: []string{"x"},
				wantBody:   "return x;",
			},
			"Two": {
				input:      "fn(x, y) { return x + y; }",
				wantParams: []string{"x", "y"},
				wantBody:   "return (x + y);",
			},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				p := parser.New(lexer.New(tc.input))
				program := p.Parse()
				checkErrors(t, p.Errors())
				is.Equal(t, 1, len(program.Statements))
				fn := program.Statements[0].(ast.ExprStmt).Expression().(ast.Function)

				params := fn.Params
				is.Equal(t, len(tc.wantParams), len(params))
				for i, want := range tc.wantParams {
					is.Equal(t, want, params[i].String())
				}
				is.Equal(t, tc.wantBody, fn.Body.String())
			})
		}
	})
	t.Run("Call Expression", func(t *testing.T) {
		t.Parallel()
		p := parser.New(lexer.New("add(1, 2 * 3, 4 + 5);"))
		program := p.Parse()
		checkErrors(t, p.Errors())
		is.Equal(t, 1, len(program.Statements))
		expr := program.Statements[0].(ast.ExprStmt)
		call := expr.Expression().(ast.CallExpr)
		is.Equal(t, "add", call.Fn.String())
		is.Equal(t, 3, len(call.Args))
		for i, want := range []string{"1", "(2 * 3)", "(4 + 5)"} {
			is.Equal(t, want, call.Args[i].String())
		}
	})
	t.Run("Strings", func(t *testing.T) {
		t.Parallel()
		p := parser.New(lexer.New(`"hey young world";`))
		program := p.Parse()
		checkErrors(t, p.Errors())
		is.Equal(t, 1, len(program.Statements))
		expr := program.Statements[0].(ast.ExprStmt)
		s := expr.Expression().(ast.String)
		is.Equal(t, "hey young world", s.String())
	})
	t.Run("Slice", func(t *testing.T) {
		t.Parallel()
		p := parser.New(lexer.New("[1, 2 * 2, 3]"))
		program := p.Parse()
		checkErrors(t, p.Errors())
		is.Equal(t, 1, len(program.Statements))
		slice := program.Statements[0].(ast.ExprStmt).Expression().(ast.Slice)
		is.Equal(t, 3, len(slice.Values()))
	})
}

func checkErrors(t *testing.T, errors []string) {
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, e := range errors {
		t.Errorf("\t%s", e)
	}
	t.FailNow()
}
