package object

import (
	"github.com/MarkyMan4/yetti/ast"
)

type FunctionObject struct {
	Args       []string
	Statements []ast.Statement
}

func (f *FunctionObject) Type() string {
	return FUNCTION_OBJ
}

func (f *FunctionObject) ToString() string {
	return "function" // TODO: implement this later
}
