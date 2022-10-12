package parser

import (
	"fmt"
	"strconv"

	"github.com/MarkyMan4/yetti/ast"
	"github.com/MarkyMan4/yetti/lexer"
	"github.com/MarkyMan4/yetti/token"
)

type (
	prefixParser func() ast.Expression
	infixParser  func(ast.Expression) ast.Expression
)

type Parser struct {
	Lex           *lexer.Lexer
	prevToken     token.Token
	curToken      token.Token
	peekToken     token.Token
	Errors        []string
	prefixParsers map[string]prefixParser
	infixParsers  map[string]infixParser
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{Lex: l}

	// load two tokens
	// prevToken will be nil, curToken and peek token will have the first two tokens
	p.nextToken()
	p.nextToken()

	// prefix parsers (e.g. an int is a prefix in an expression)
	p.prefixParsers = map[string]prefixParser{
		token.INT:     p.parseIntegerLiteral,
		token.FLOAT:   p.parseFloatLiteral,
		token.BOOLEAN: p.parseBooleanLiteral,
		token.STRING:  p.parseStringLiteral,
		token.IDENT:   p.parseIdent,
		token.LBRACK:  p.parseArray,
	}

	// infix parsers (e.g. +, -, *, /)
	p.infixParsers = map[string]infixParser{
		token.PLUS:   p.parseInfixExpression,
		token.PLUSEQ: p.parseInfixExpression,
		token.MINUS:  p.parseInfixExpression,
		token.MINEQ:  p.parseInfixExpression,
		token.MULT:   p.parseInfixExpression,
		token.MULTEQ: p.parseInfixExpression,
		token.DIVIDE: p.parseInfixExpression,
		token.DIVEQ:  p.parseInfixExpression,
		token.LT:     p.parseInfixExpression,
		token.LTE:    p.parseInfixExpression,
		token.EQ:     p.parseInfixExpression,
		token.GT:     p.parseInfixExpression,
		token.GTE:    p.parseInfixExpression,
		token.DOT:    p.parseObjFuncExpression,
		token.LBRACK: p.parseIndexExpression,
	}

	return p
}

func (p *Parser) nextToken() {
	p.prevToken = p.curToken
	p.curToken = p.peekToken
	p.peekToken = p.Lex.NextToken()
}

func (p *Parser) Parse() *ast.Program {
	prog := &ast.Program{Statements: []ast.Statement{}}

	for p.curToken.Type != token.EOF {
		prog.Statements = append(prog.Statements, p.parseStmt())
		p.nextToken()

		if p.curToken.Type == token.SEMI {
			p.nextToken()
		}
	}

	return prog
}

func (p *Parser) parseStmt() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStmt()
	case token.WHILE:
		return p.parseWhileStmt()
	case token.IF:
		return p.parseIfStmt()
	case token.FUN:
		return p.parseFunctionDef()
	case token.RETURN:
		return p.parseReturnStmt()
	case token.IDENT:
		// determine how the ident is being used based on the next token
		switch p.peekToken.Literal {
		case token.LPAREN:
			return p.parseFunctionCall()
		default:
			return p.parseAssignStmt()
		}
	default:
		return nil
	}
}

func (p *Parser) parseVarStmt() ast.Statement {
	if !p.expectNextToken(token.IDENT) {
		return nil
	}

	p.nextToken()
	stmt := &ast.VarStatement{Identifier: p.curToken.Literal}

	if !p.expectNextToken(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	p.nextToken()
	stmt.Value = p.parseExpression()

	// TODO: figure out where cursor should leave off after parsing function call (on semi colon or right paren)
	if p.curToken.Type != token.SEMI && !p.expectNextToken(token.SEMI) {
		return nil
	}

	if p.curToken.Type != token.SEMI {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) expectNextToken(tokType string) bool {
	if p.peekToken.Type == tokType {
		return true
	}

	errMsg := fmt.Sprintf("Expected token %s to be %s, but got %s", p.peekToken, tokType, p.peekToken.Type)
	p.Errors = append(p.Errors, errMsg)

	return false
}

func (p *Parser) parseExpression() ast.Expression {
	prefix := p.prefixParsers[p.curToken.Type]
	if prefix == nil {
		// TODO: do some kind of error here
		return nil
	}

	left := prefix()
	for p.peekToken.Type != token.SEMI {
		infix := p.infixParsers[p.peekToken.Type]
		if infix == nil {
			return left
		}

		p.nextToken()
		left = infix(left)
	}

	return left
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	intLit := &ast.IntegerLiteral{}
	val, err := strconv.ParseInt(p.curToken.Literal, 10, 64)

	if err != nil {
		errMsg := fmt.Sprintf("could not parse %s as type integer", p.curToken.Literal)
		p.Errors = append(p.Errors, errMsg)
		return nil
	}

	intLit.Value = val

	return intLit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	floatLit := &ast.FloatLiteral{}
	val, err := strconv.ParseFloat(p.curToken.Literal, 64)

	if err != nil {
		errMsg := fmt.Sprintf("could not parse %s as type float", p.curToken.Literal)
		p.Errors = append(p.Errors, errMsg)
		return nil
	}

	floatLit.Value = val

	return floatLit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	boolLit := &ast.BooleanLiteral{}
	val, err := strconv.ParseBool(p.curToken.Literal)

	if err != nil {
		errMsg := fmt.Sprintf("could not parse %s as type boolean", p.curToken.Literal)
		p.Errors = append(p.Errors, errMsg)
		return nil
	}

	boolLit.Value = val

	return boolLit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Value: p.curToken.Literal}
}

// handles parsing variables and function calls
func (p *Parser) parseIdent() ast.Expression {
	if p.peekToken.Type == token.LPAREN {
		res := p.parseFunctionCall().(ast.Expression)
		return res
	}

	return &ast.IdentifierExpression{Value: p.curToken.Literal}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Op:   p.curToken.Literal,
		Left: left,
	}

	p.nextToken()
	expr.Right = p.parseExpression()

	return expr
}

func (p *Parser) parseWhileStmt() ast.Statement {
	whileStmt := &ast.WhileStatement{Statements: []ast.Statement{}}

	if !p.expectNextToken(token.LPAREN) {
		// may have to skip to end of line?
		return nil
	}

	p.nextToken()
	p.nextToken()
	whileStmt.Condition = p.parseExpression()
	if !p.expectNextToken(token.RPAREN) {
		return nil
	}

	p.nextToken()
	if !p.expectNextToken((token.LBRACE)) {
		return nil
	}

	p.nextToken()
	p.nextToken()
	// TODO: handle missing RBRACE
	for p.curToken.Type != token.RBRACE {
		whileStmt.Statements = append(whileStmt.Statements, p.parseStmt())
		p.nextToken()
	}

	return whileStmt
}

func (p *Parser) parseIfStmt() ast.Statement {
	ifStmt := &ast.IfStatement{Statements: []ast.Statement{}}

	if !p.expectNextToken(token.LPAREN) {
		// may have to skip to end of line?
		return nil
	}

	p.nextToken()
	p.nextToken()
	ifStmt.Condition = p.parseExpression()
	if !p.expectNextToken(token.RPAREN) {
		return nil
	}

	p.nextToken()
	if !p.expectNextToken((token.LBRACE)) {
		return nil
	}

	p.nextToken()
	p.nextToken()
	for p.curToken.Type != token.RBRACE {
		ifStmt.Statements = append(ifStmt.Statements, p.parseStmt())
		p.nextToken()
	}

	return ifStmt
}

func (p *Parser) parseReturnStmt() ast.Statement {
	p.nextToken()
	returnStmt := &ast.ReturnStatement{ReturnVal: p.parseExpression()}

	if !p.expectNextToken(token.SEMI) {
		return nil
	}

	p.nextToken()

	return returnStmt
}

func (p *Parser) parseAssignStmt() ast.Statement {
	assignStmt := &ast.AssignStatement{Identifier: p.curToken.Literal}

	// TODO: don't use expectNextToken() here since it will cause unwanted errors
	if !p.expectNextToken(token.ASSIGN) &&
		!p.expectNextToken(token.PLUSEQ) &&
		!p.expectNextToken(token.MINEQ) &&
		!p.expectNextToken(token.MULTEQ) &&
		!p.expectNextToken(token.DIVEQ) {
		return nil
	}

	p.nextToken()
	assignStmt.AssignOp = p.curToken.Literal
	p.nextToken()

	assignStmt.Value = p.parseExpression()

	if !p.expectNextToken(token.SEMI) {
		return nil
	}

	p.nextToken()

	return assignStmt
}

// TODO: handle function with no args
func (p *Parser) parseFunctionDef() ast.Statement {
	if !p.expectNextToken(token.IDENT) {
		return nil
	}

	p.nextToken()

	funcDef := &ast.FunctionDef{
		Name:       p.curToken.Literal,
		Args:       []string{},
		Statements: []ast.Statement{},
	}

	if !p.expectNextToken(token.LPAREN) {
		return nil
	}

	p.nextToken()
	p.nextToken()

	for p.curToken.Type == token.IDENT {
		funcDef.Args = append(funcDef.Args, p.curToken.Literal)

		if p.peekToken.Type != token.RPAREN && p.peekToken.Type != token.COM {
			errMsg := fmt.Sprintf("Unexpected token %s. Expected %s or %s", p.peekToken, token.RPAREN, token.COM)
			p.Errors = append(p.Errors, errMsg)

			return nil
		}

		p.nextToken()
		p.nextToken()
	}

	// in case there were no arguments, the loop was skipped
	if p.curToken.Type == token.RPAREN {
		p.nextToken()
	}

	if p.curToken.Type != token.LBRACE {
		errMsg := fmt.Sprintf("Unexpected token %s. Expected %s", p.curToken.Literal, token.LBRACE)
		p.Errors = append(p.Errors, errMsg)

		return nil
	}

	p.nextToken()

	// TODO: handle missing rbrace
	for p.curToken.Type != token.RBRACE {
		funcDef.Statements = append(funcDef.Statements, p.parseStmt())
		p.nextToken()

		if p.curToken.Type == token.SEMI {
			p.nextToken()
		}
	}

	return funcDef
}

func (p *Parser) parseFunctionCall() ast.Statement {
	funcCall := &ast.FunctionCall{Name: p.curToken.Literal, Args: []ast.Expression{}}
	if !p.expectNextToken(token.LPAREN) {
		return nil
	}

	p.nextToken()
	p.nextToken()

	for p.curToken.Type != token.RPAREN && p.curToken.Type != token.SEMI && p.curToken.Type != token.EOF {
		funcCall.Args = append(funcCall.Args, p.parseExpression())

		if p.peekToken.Type != token.RPAREN && p.peekToken.Type != token.COM {
			errMsg := fmt.Sprintf("Unexpected token %s. Expected %s or %s", p.peekToken, token.RPAREN, token.COM)
			p.Errors = append(p.Errors, errMsg)

			return nil
		}

		p.nextToken()

		if p.curToken.Type != token.RPAREN {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.EOF {
		errMsg := fmt.Sprintf("Unexpected token %s.", p.curToken.Literal)
		p.Errors = append(p.Errors, errMsg)

		return nil
	}

	return funcCall
}

// this is just an infix expression where left is an object, op is '.' and right is the function call for that object
func (p *Parser) parseObjFuncExpression(obj ast.Expression) ast.Expression {
	fnCall := &ast.ObjectFunctionExpression{Object: obj}
	p.nextToken()
	fnCall.Function = p.parseIdent()

	return fnCall
}

func (p *Parser) parseArray() ast.Expression {
	p.nextToken()
	arr := &ast.ArrayExpression{Items: []ast.Expression{}}

	for p.curToken.Type != token.RBRACK && p.curToken.Type != token.SEMI && p.curToken.Type != token.EOF {
		arr.Items = append(arr.Items, p.parseExpression())

		if p.peekToken.Type != token.RBRACK && p.peekToken.Type != token.COM {
			errMsg := fmt.Sprintf("Unexpected token %s. Expected %s or %s", p.peekToken, token.RBRACK, token.COM)
			p.Errors = append(p.Errors, errMsg)

			return nil
		}

		p.nextToken()

		if p.curToken.Type != token.RBRACK {
			p.nextToken()
		}
	}

	if p.curToken.Type == token.EOF {
		errMsg := fmt.Sprintf("Unexpected token %s.", p.curToken.Literal)
		p.Errors = append(p.Errors, errMsg)

		return nil
	}

	return arr
}

func (p *Parser) parseIndexExpression(arr ast.Expression) ast.Expression {
	idxExpr := &ast.ArrayIndexExpression{Arr: arr}
	p.nextToken()
	idxExpr.Index = p.parseExpression()

	if p.curToken.Type == token.RBRACK {
		errMsg := fmt.Sprintf("Unexpected token %s.", p.curToken.Literal)
		p.Errors = append(p.Errors, errMsg)

		return nil
	}

	p.nextToken()

	return idxExpr
}
