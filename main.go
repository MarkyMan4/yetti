package main

import (
	"fmt"
	"os"

	"github.com/MarkyMan4/yetti/evaluator"
	"github.com/MarkyMan4/yetti/lexer"
	"github.com/MarkyMan4/yetti/object"
	"github.com/MarkyMan4/yetti/parser"
)

func readFile(filename string) string {
	file, err := os.ReadFile(filename)

	if err != nil {
		fmt.Printf("error reading file %s\n", filename)
		os.Exit(1)
	}

	return string(file)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("you must provide a filename")
		os.Exit(1)
	}

	text := readFile(os.Args[1])

	env := object.NewEnvironment()
	l := lexer.NewLexer(text)
	p := parser.NewParser(l)
	prog := p.Parse()

	// evaluate each statement in the program
	for i := range prog.Statements {
		evaluator.Eval(prog.Statements[i], env)
	}

	// print out the state of the program
	// for k, v := range env.GetEnvMap() {
	// 	fmt.Printf("%s: %s\n", k, v.ToString())
	// }
}
