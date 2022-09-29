package parser

import (
	"fmt"
	"testing"

	"github.com/MarkyMan4/yetti/ast"
	"github.com/MarkyMan4/yetti/lexer"
)

// func TestParse(t *testing.T) {
// 	l := lexer.NewLexer("var x = 1; while(x < 5) {x += 1;}")
// 	// l := lexer.NewLexer("var x = 1; someFunc(x, 2);")
// 	// l := lexer.NewLexer("var x = 1;")
// 	p := NewParser(l)
// 	prog := p.Parse()

// 	// stmt := prog.Statements[1].(*ast.WhileStatement)

// 	// fmt.Println(stmt.Name)

// 	for i := range prog.Statements {
// 		fmt.Println(prog.Statements[i])
// 	}

// 	// fmt.Println(stmt.ToString())

// 	// for i := range stmt.Statements {
// 	// 	fmt.Println(stmt.Statements[i])
// 	// }

// 	// tok := l.NextToken()

// 	// for tok.Type != token.EOF {
// 	// 	fmt.Println(tok)
// 	// 	tok = l.NextToken()
// 	// }
// }

// func TestParseIf(t *testing.T) {
// 	l := lexer.NewLexer("var x = 1; if(x < 5) {x += 1;}")
// 	// l := lexer.NewLexer("var x = 1; someFunc(x, 2);")
// 	// l := lexer.NewLexer("var x = 1;")
// 	p := NewParser(l)
// 	prog := p.Parse()

// 	stmt := prog.Statements[1].(*ast.IfStatement)

// 	fmt.Println(stmt.ToString())

// 	// for i := range prog.Statements {
// 	// 	fmt.Println(prog.Statements[i])
// 	// }

// 	// fmt.Println(stmt.ToString())

// 	// for i := range stmt.Statements {
// 	// 	fmt.Println(stmt.Statements[i])
// 	// }
// }

func TestParseFuncDef(t *testing.T) {
	l := lexer.NewLexer("fun test(x, y) {var a = x; b = y;} var x = test(1, 2);")
	p := NewParser(l)
	prog := p.Parse()

	// stmt := prog.Statements[0].(*ast.FunctionDef)

	// fmt.Println(stmt.ToString())

	// stmt1 := prog.Statements[1].(*ast.VarStatement)

	// fmt.Println(stmt1.ToString())

	for i := range prog.Statements {
		fmt.Println(prog.Statements[i])
	}

	// fmt.Println(stmt.ToString())

	// for i := range stmt.Statements {
	// 	fmt.Println(stmt.Statements[i])
	// }
}

func TestParseString(t *testing.T) {
	fmt.Println("------ test string parse -------")
	l := lexer.NewLexer("print(\"hello\");")
	p := NewParser(l)
	prog := p.Parse()

	stmt := prog.Statements[0].(*ast.FunctionCall)

	fmt.Println(stmt.ToString())
}

func TestParseObjectFunctionCall(t *testing.T) {
	fmt.Println("------ test object function call -------")
	l := lexer.NewLexer("var s = \"hello\"; var x = s.substring(1, 2);")
	p := NewParser(l)
	prog := p.Parse()

	stmt := prog.Statements[1].(*ast.VarStatement)

	fmt.Println(stmt.ToString())
}

func TestParseArray(t *testing.T) {
	fmt.Println("------ test parsing array -------")
	l := lexer.NewLexer("var xs = [1,2,3,4];")
	p := NewParser(l)
	prog := p.Parse()

	stmt := prog.Statements[0].(*ast.VarStatement)

	fmt.Println(stmt.ToString())
}

func TestParseArrayIndex(t *testing.T) {
	fmt.Println("------ test parsing array index -------")
	l := lexer.NewLexer("var xs = [1,2,3,4]; var i = xs[2]; var y = 3.45;")
	p := NewParser(l)
	prog := p.Parse()

	for i := range prog.Statements {
		stmt := prog.Statements[i].(*ast.VarStatement)
		fmt.Println(stmt.ToString())
	}
}
