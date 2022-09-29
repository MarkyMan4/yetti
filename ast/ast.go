package ast

import "fmt"

type Node interface {
	ToString() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// expressions
type IntegerLiteral struct {
	Value int64
}

func (i *IntegerLiteral) ToString() string {
	return fmt.Sprint(i.Value)
}

func (i *IntegerLiteral) expressionNode() {}

type FloatLiteral struct {
	Value float64
}

func (i *FloatLiteral) ToString() string {
	return fmt.Sprint(i.Value)
}

func (i *FloatLiteral) expressionNode() {}

type StringLiteral struct {
	Value string
}

func (i *StringLiteral) ToString() string {
	return i.Value
}

func (b *StringLiteral) expressionNode() {}

type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) ToString() string {
	return fmt.Sprint(b.Value)
}

func (b *BooleanLiteral) expressionNode() {}

type InfixExpression struct {
	Left  Expression
	Op    string
	Right Expression
}

func (i *InfixExpression) ToString() string {
	// fix this later, just returning empty string for now
	return i.Left.ToString() + i.Op + i.Right.ToString()
}

func (i *InfixExpression) expressionNode() {}

type IdentifierExpression struct {
	Value string // name of identifier
}

func (i *IdentifierExpression) ToString() string {
	return i.Value
}

func (i *IdentifierExpression) expressionNode() {}

// function calls on an object, e.g. var a = "hello"; var b = s.substring(1, 3);
type ObjectFunctionExpression struct {
	Object   Expression
	Function Expression
}

func (o *ObjectFunctionExpression) ToString() string {
	return fmt.Sprintf("%s.%s", o.Object.ToString(), o.Function.ToString())
}

func (o *ObjectFunctionExpression) expressionNode() {}

type ArrayExpression struct {
	Items []Expression
}

func (a *ArrayExpression) ToString() string {
	arrStr := ""

	for i := range a.Items {
		arrStr += a.Items[i].ToString()
		if i < len(a.Items)-1 {
			arrStr += ","
		}
	}

	return fmt.Sprintf("[%s]", arrStr)
}

func (a *ArrayExpression) expressionNode() {}

// e.g. var arr = [1,2,3]; var i = arr[0];
type ArrayIndexExpression struct {
	Arr   Expression
	Index Expression
}

func (a *ArrayIndexExpression) ToString() string {
	return fmt.Sprintf("%s[%s]", a.Arr.ToString(), a.Index.ToString())
}

func (a *ArrayIndexExpression) expressionNode() {}

// statements
type VarStatement struct {
	Identifier string
	Value      Expression
}

func (l *VarStatement) ToString() string {
	return fmt.Sprintf("var %s = %s;", l.Identifier, l.Value.ToString())
}

func (l *VarStatement) statementNode() {}

type AssignStatement struct {
	Identifier string
	AssignOp   string // assignment operators are =, +=, -=, *=, /=
	Value      Expression
}

func (a *AssignStatement) ToString() string {
	return fmt.Sprintf("%s %s %s;", a.Identifier, a.AssignOp, a.Value.ToString())
}

func (a *AssignStatement) statementNode() {}

// function invocation
type FunctionCall struct {
	Name string
	Args []Expression
}

func (fc *FunctionCall) ToString() string {
	argsStr := ""

	for i := range fc.Args {
		argsStr += fc.Args[i].ToString()
		if i < len(fc.Args)-1 {
			argsStr += ", "
		}
	}

	funcCallStr := fmt.Sprintf("%s(%s)", fc.Name, argsStr)

	return funcCallStr
}

// function call can be used as both a statement and expression
func (fc *FunctionCall) expressionNode() {}
func (fc *FunctionCall) statementNode()  {}

// while loop
type WhileStatement struct {
	Condition  Expression
	Statements []Statement
}

func (ws *WhileStatement) ToString() string {
	whileAsStr := fmt.Sprintf("while(%s) {", ws.Condition.ToString())

	for i := range ws.Statements {
		whileAsStr += ws.Statements[i].ToString() + " "
	}

	whileAsStr += "}"

	return whileAsStr
}

func (ws *WhileStatement) statementNode() {}

// if statement
type IfStatement struct {
	Condition  Expression
	Statements []Statement
}

func (is *IfStatement) ToString() string {
	ifAsStr := fmt.Sprintf("if(%s) { ", is.Condition.ToString())

	for i := range is.Statements {
		ifAsStr += is.Statements[i].ToString() + " "
	}

	ifAsStr += "}"

	return ifAsStr
}

func (is *IfStatement) statementNode() {}

// function definition
type FunctionDef struct {
	Name       string
	Args       []string // list of identifiers
	Statements []Statement
}

func (fd *FunctionDef) ToString() string {
	funcStr := fmt.Sprintf("def %s(", fd.Name)

	for i := range fd.Args {
		funcStr += fd.Args[i]
		if i < len(fd.Args)-1 {
			funcStr += ","
		} else {
			funcStr += ")"
		}
	}

	funcStr += " { "

	for i := range fd.Statements {
		funcStr += fd.Statements[i].ToString() + " "
	}

	funcStr += "}"

	return funcStr
}

func (fd *FunctionDef) statementNode() {}

// return statement
type ReturnStatement struct {
	ReturnVal Expression
}

func (rs *ReturnStatement) ToString() string {
	return fmt.Sprintf("return %s;", rs.ReturnVal.ToString())
}

func (rs *ReturnStatement) statementNode() {}

// program is a list of statements
type Program struct {
	Statements []Statement
}

func (p *Program) ToString() string {
	return ""
}
