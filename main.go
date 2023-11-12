package main

import (
	"mmm/repl"
	"os"
)

func main() {
	os.Stdout.Write([]byte("Mmm monkey\n"))
	repl.Start(os.Stdin, os.Stdout)
}
