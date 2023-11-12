package repl

import (
	"bufio"
	"fmt"
	"io"
	"mmm/entity"
	"mmm/eval"
	"mmm/lexer"
	"mmm/parser"
	"mmm/token"
)

const prompt = ">> "

func Start(in io.Reader, out io.Writer) {
	s := bufio.NewScanner(in)
	env := entity.NewEnv()
	for s.Scan() {
		fmt.Fprint(out, prompt)
		p := parser.New(lexer.New((s.Text())))
		prg := p.Parse()
		if len(p.Errors()) != 0 {
			for _, e := range p.Errors() {
				fmt.Fprint(out, "\t"+e+"\n")
			}
			continue
		}
		if e := eval.Eval(prg, env); e != nil {
			fmt.Fprintf(out, "%+v\n", e.Inspect())
		}
	}
}

func StartParser(in io.Reader, out io.Writer) {
	s := bufio.NewScanner(in)
	for s.Scan() {
		fmt.Fprint(out, prompt)
		p := parser.New(lexer.New((s.Text())))
		prg := p.Parse()
		if len(p.Errors()) != 0 {
			for _, e := range p.Errors() {
				fmt.Fprint(out, "\t"+e+"\n")
			}
			continue
		}
		fmt.Fprintf(out, "%+v\n", prg.String())
	}
}

func StartLexer(in io.Reader, out io.Writer) {
	s := bufio.NewScanner(in)
	for s.Scan() {
		fmt.Fprint(out, prompt)
		l := lexer.New(s.Text())
		for t := l.NextToken(); t.Type() != token.TypeEOF; t = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", t)
		}
	}
}
