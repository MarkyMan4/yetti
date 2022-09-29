package main

import (
	"fmt"
	"testing"

	"github.com/MarkyMan4/yetti/ast"
)

func getObj() ast.Statement {
	return &ast.ReturnStatement{ReturnVal: &ast.IntegerLiteral{Value: 4}}
}

func TestMain(t *testing.T) {
	// m := map[string]string{}

	// if res, ok := m["test"]; !ok {
	// 	fmt.Println("nil")
	// } else {
	// 	fmt.Println(res)
	// }

	obj := getObj()

	switch obj.(type) {
	case *ast.ReturnStatement:
		fmt.Println("yea")
	default:
		fmt.Println("no")
	}
}
