package eval_test

import (
	"mmm/entity"
	"mmm/eval"
	"mmm/is"
	"mmm/lexer"
	"mmm/parser"
	"testing"
)

func TestEval(t *testing.T) {
	t.Parallel()
	t.Run("Integers", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input string
			want  int64
		}{
			"One digit":           {input: "5", want: 5},
			"Two digits":          {input: "10", want: 10},
			"Negative One digit":  {input: "-5", want: -5},
			"Negative Two digits": {input: "-10", want: -10},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				is.Equal(t, tc.want, ent.(entity.Int).Value)
			})
		}
	})
	t.Run("Infix Bool Expressions", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input string
			want  bool
		}{
			"True":                    {input: "true", want: true},
			"False":                   {input: "false", want: false},
			"True is true":            {input: "true == true", want: true},
			"False is false":          {input: "false == false", want: true},
			"LT":                      {input: "1 < 2", want: true},
			"GT":                      {input: "1 > 2", want: false},
			"LT on EQ":                {input: "1 < 1", want: false},
			"GT on EQ":                {input: "1 > 1", want: false},
			"EQ":                      {input: "1 == 1", want: true},
			"Not EQ":                  {input: "1 != 1", want: false},
			"EQ Differnt Values":      {input: "1 == 2", want: false},
			"Not EQ Different Values": {input: "1 != 2", want: true},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				is.Equal(t, tc.want, ent.(entity.Bool).Value)
			})
		}
	})
	t.Run("Bang Operator", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input string
			want  bool
		}{
			"Not True":         {input: "!true", want: false},
			"Not False":        {input: "!false", want: true},
			"Not Five":         {input: "!5", want: false},
			"Double Not True":  {input: "!!true", want: true},
			"Double Not False": {input: "!!false", want: false},
			"Double Not Five":  {input: "!!5", want: true},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				is.Equal(t, tc.want, ent.(entity.Bool).Value)
			})
		}
	})
	t.Run("Infix Integer Expressions", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input string
			want  int64
		}{
			"Add":          {input: "5 + 5", want: 10},
			"Subtract":     {input: "5 - 5", want: 0},
			"Multiply":     {input: "5 * 5", want: 25},
			"Divide":       {input: "5 / 5", want: 1},
			"Kitchen Sink": {input: "5 * (5 + 5) - 55 / 5", want: 39},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				is.Equal(t, tc.want, ent.(entity.Int).Value)
			})
		}
	})
	t.Run("If-Else Expressions", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input string
			want  int64
		}{
			"Literal true":  {input: "if (true) { 10 }", want: 10},
			"Literal false": {input: "if (false) { 10 }", want: -1},
			"Truthy int":    {input: "if (1) { 10 }", want: 10},
			"LT":            {input: "if (1 < 2) { 10 }", want: 10},
			"GT":            {input: "if (1 > 2) { 10 }", want: -1},
			"LT with else":  {input: "if (1 < 2) { 10 } else { 20 }", want: 10},
			"GT with else":  {input: "if (1 > 2) { 10 } else { 20 }", want: 20},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				if tc.want == -1 {
					is.Equal(t, entity.TypeNull, ent.Type())
					return
				}
				is.Equal(t, tc.want, ent.(entity.Int).Value)
			})
		}
	})
	t.Run("Return statements", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input string
			want  int64
		}{
			"Single statement":                  {input: "return 10;", want: 10},
			"Statement after return":            {input: "return 10; 9;", want: 10},
			"Statement before and after return": {input: "9;return 10;9;", want: 10},
			"Statement as expression":           {input: "return 2 * 5", want: 10},
			"No immediate return": {input: `
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
				return 1;
			}`, want: 10},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				is.Equal(t, tc.want, ent.(entity.Int).Value)
			})
		}
	})
	t.Run("Error Handling", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input string
			want  string
		}{
			"Add Int to Bool": {input: "5 + true;",
			want: "type mismatch: Int + Bool"},
			"Add Int to Bool; Statement": {input: "5 + true; 5;",
			want: "type mismatch: Int + Bool"},
			"Minus Bool": {input: "-true;", want: "unknown operator: -Bool"},
			"Add Bool to Bool": {input: "true + true;",
			want: "unknown operator: Bool + Bool"},
			"Add Bool to Bool in If": {input: "if (10 > 1) { true + false; }",
			want: "unknown operator: Bool + Bool"},
			"Add Bool to Bool in Nested If": {input: `
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}`,
			want: "unknown operator: Bool + Bool"},
			"No var foo in environment": {input: "foo;",
			want: "identifier not found: foo",},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				is.Equal(t, tc.want, ent.(entity.Error).Message)
			})
		}
	})
	t.Run("Let statements", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input string
			want  int64
		}{
			"Int Literal":          {input: "let a = 5; a;", want: 5},
			"Expression":          {input: "let a = 5 * 5; a;", want: 25},
			"Two Lets":          {input: "let a = 5; let b = a; b;", want: 5},
			"Kitchen Sink": {input: "let a=5; let b=a; let c=a+b+5; c;", want: 15},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				is.Equal(t, tc.want, ent.(entity.Int).Value)
			})
		}
	})
	t.Run("Function Expression", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input string
			want  string
		}{
			"Simple": {input: "fn(x) { return x + 2; };", want: "return (x + 2);" },
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				is.Equal(t, tc.want, ent.(entity.Fn).Body.String())
			})
		}
	})
	t.Run("Call Function", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input string
			want  int64
		}{
			"Simple": {
				input: "let id = fn(x) { return x; }; id(5);",
				want: 5,
			},
			"Double param": {
				input: "let dbl = fn(x) { return x * 2; }; dbl(5);",
				want: 10,
			},
			"Two params": {
				input: "let add = fn(x, y) { return x + y; }; add(5, 5);",
				want: 10,
			},
			"Call in Call": {
				input: "let add = fn(x, y) { return x + y; }; add(5 + 5, add(5, 5));",
				want: 20,
			},
			"Anonymous call": {
				input: "fn(x, y) { return x + y; }(5, 5);",
				want: 10,
			},
			"Closures": {
				input: `
let newAdder = fn(x) {
	return fn(y) { return x + y; };
};
let addTwo = newAdder(2);
addTwo(2);`,
				want: 4,
			},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				is.Equal(t, tc.want, ent.(entity.Int).Value)
			})
		}
	})
	t.Run("String", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input string
			want  string
		}{
			"First Words": {input:`"Hey Young Wurld!"`, want: "Hey Young Wurld!" },
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				is.Equal(t, tc.want, ent.(entity.String).Value)
			})
		}
	})
	t.Run("String Concatenation", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input string
			want  string
		}{
			"First Words": {
				input:`"Hey" + " Young Wurld!"`, want: "Hey Young Wurld!" },
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				is.Equal(t, tc.want, ent.(entity.String).Value)
			})
		}
	})
	t.Run("Builtin Len", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input string
			want  int64
		}{
			"Len: Empty": { input:`len("")`, want: 0},
			"Len: Long string": { input:`len("Hey Yung Wurld!")`, want: 15},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				is.Equal(t, tc.want, ent.(entity.Int).Value)
			})
		}
		for name, tc := range map[string]struct {
			input string
			want  string
		}{
			"Len: Bad type": { input:`len(1)`,
			want: "ERROR: argument to `len` not supported, got Int"},
			"Len: Two params": { input:`len("Hey", " Yung Wurld!")`,
			want: "ERROR: len only accepts one argument."},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				is.Equal(t, tc.want, ent.(entity.Error).Inspect())
			})
		}
	})
	t.Run("Slices", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input string
			want  []int64
		}{
			"Empty": { input: "[]", want: []int64{}},
			"One Element": { input: "[1]", want: []int64{1}},
			"Many Elements": { input: "[1,1+1,3]", want: []int64{1,2,3}},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				for i, want := range tc.want {
				is.Equal(t, want, ent.(entity.Slice).Values[i].(entity.Int).Value)
				}
			})
		}
	})
	t.Run("Indexes", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct {
			input string
			want  int64
		}{
			"First": { input: "[1,2,3][0]", want: 1},
			"As Expression": { input: "[1,2,3][1+0]", want: 2},
			"As Env": { input: "let i=0;[1,2,3][i]", want: 1},
			"As Literal": { input: "let i=0;[1,2,3][i] + 1", want: 2},
			"No Negative": { input: "[1,2,3][-1]", want: -1},
			"OOB": { input: "[1,2,3][4]", want: -1},
		} {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				ent := setup(tc.input)
				if tc.want == -1 {
				is.Equal(t, entity.TypeNull, ent.(entity.Null).Type())
					return
				}
			is.Equal(t, tc.want, ent.(entity.Int).Value)
			})
		}
	})
}

func setup(input string) entity.E {
	return eval.Eval(parser.New(lexer.New(input)).Parse(), entity.NewEnv())
}
