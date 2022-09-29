package evaluator

import (
	"fmt"
	"os"

	"github.com/MarkyMan4/yetti/ast"
	"github.com/MarkyMan4/yetti/object"
	"github.com/MarkyMan4/yetti/stdlib"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.IntegerObject{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.FloatObject{Value: node.Value}
	case *ast.StringLiteral:
		return &object.StringObject{Value: node.Value}
	case *ast.BooleanLiteral:
		return &object.BooleanObject{Value: node.Value}
	case *ast.ArrayExpression:
		return evalArrayExpression(node.Items, env)
	case *ast.ArrayIndexExpression:
		return evalArrayIndexExpression(node, env)
	case *ast.IdentifierExpression:
		obj, ok := env.Get(node.Value)

		if !ok {
			fmt.Printf("identifier %s is not defined\n", node.Value)
			os.Exit(1)
		}

		return obj
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalInfixExpression(node.Op, left, right)
	case *ast.VarStatement:
		val := Eval(node.Value, env)
		env.Set(node.Identifier, val, true)
	case *ast.AssignStatement:
		obj, ok := env.Get(node.Identifier)

		if !ok {
			fmt.Printf("variable %s has not been declared\n", node.Identifier)
			os.Exit(1)
		}

		left := obj
		right := Eval(node.Value, env)
		val := evalAssignStatement(node.AssignOp, left, right)
		env.Set(node.Identifier, val, false)
	case *ast.IfStatement:
		condResult := Eval(node.Condition, env)
		if condResult.Type() != object.BOOLEAN_OBJ {
			fmt.Println("condition must return a boolean")
			os.Exit(1)
		}

		// if the condition is still true, run all statements and evaluate the loop again
		if condResult.(*object.BooleanObject).Value {
			for i := range node.Statements {
				Eval(node.Statements[i], env)
			}
		}
	case *ast.WhileStatement:
		condResult := Eval(node.Condition, env)
		if condResult.Type() != object.BOOLEAN_OBJ {
			fmt.Println("condition must return a boolean")
			os.Exit(1)
		}

		// if the condition is still true, run all statements and evaluate the loop again
		if condResult.(*object.BooleanObject).Value {
			for i := range node.Statements {
				Eval(node.Statements[i], env)
			}

			Eval(node, env)
		}
	case *ast.FunctionDef:
		env.Set(node.Name, &object.FunctionObject{Args: node.Args, Statements: node.Statements}, true)
	case *ast.FunctionCall:
		return evalFunctionCall(node, env)
	case *ast.ReturnStatement:
		return Eval(node.ReturnVal, env)
	case *ast.ObjectFunctionExpression:
		return evalObjFunCall(node, env)
	}

	return nil
}

func evalAssignStatement(assignOp string, left object.Object, right object.Object) object.Object {
	switch assignOp {
	case "+=":
		return evalInfixExpression("+", left, right)
	case "-=":
		return evalInfixExpression("-", left, right)
	case "*=":
		return evalInfixExpression("*", left, right)
	case "/=":
		return evalInfixExpression("/", left, right)
	case "=":
		// default is a normal assignment, so just return the right hand side
		return right
	default:
		fmt.Printf("unknown operator %s\n", assignOp)
		os.Exit(1)
	}

	return nil
}

func evalInfixExpression(op string, left object.Object, right object.Object) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return evalIntegerInfixExpression(op, left, right)
	} else if left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ {
		return evalIntegerFloatInfixExpression(op, left, right)
	} else if left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ {
		return evalFloatIntegerInfixExpression(op, left, right)
	} else if left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ {
		return evalFloatInfixExpression(op, left, right)
	}

	return nil
}

func evalIntegerInfixExpression(op string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.IntegerObject).Value
	rightVal := right.(*object.IntegerObject).Value

	switch op {
	case "+":
		return &object.IntegerObject{Value: leftVal + rightVal}
	case "-":
		return &object.IntegerObject{Value: leftVal - rightVal}
	case "*":
		return &object.IntegerObject{Value: leftVal * rightVal}
	case "/":
		// for dividing integers, convert them to floats so we get a float in return
		return &object.FloatObject{Value: float64(leftVal) / float64(rightVal)}
	case "<":
		return &object.BooleanObject{Value: leftVal < rightVal}
	case "<=":
		return &object.BooleanObject{Value: leftVal <= rightVal}
	case "==":
		return &object.BooleanObject{Value: leftVal == rightVal}
	case ">":
		return &object.BooleanObject{Value: leftVal > rightVal}
	case ">=":
		return &object.BooleanObject{Value: leftVal >= rightVal}
	default:
		return &object.ErrorObject{Message: fmt.Sprintf("unsupported operator '%s' for types %s, %s", op, left.Type(), right.Type())}
	}
}

func evalFloatIntegerInfixExpression(op string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.FloatObject).Value
	rightVal := float64(right.(*object.IntegerObject).Value)

	switch op {
	case "+":
		return &object.FloatObject{Value: leftVal + rightVal}
	case "-":
		return &object.FloatObject{Value: leftVal - rightVal}
	case "*":
		return &object.FloatObject{Value: leftVal * rightVal}
	case "/":
		return &object.FloatObject{Value: leftVal / rightVal}
	case "<":
		return &object.BooleanObject{Value: leftVal < rightVal}
	case "<=":
		return &object.BooleanObject{Value: leftVal <= rightVal}
	case "==":
		return &object.BooleanObject{Value: leftVal == rightVal}
	case ">":
		return &object.BooleanObject{Value: leftVal > rightVal}
	case ">=":
		return &object.BooleanObject{Value: leftVal >= rightVal}
	default:
		return &object.ErrorObject{Message: fmt.Sprintf("unsupported operator '%s' for types %s, %s", op, left.Type(), right.Type())}
	}
}

func evalIntegerFloatInfixExpression(op string, left object.Object, right object.Object) object.Object {
	leftVal := float64(left.(*object.IntegerObject).Value)
	rightVal := right.(*object.FloatObject).Value

	switch op {
	case "+":
		return &object.FloatObject{Value: leftVal + rightVal}
	case "-":
		return &object.FloatObject{Value: leftVal - rightVal}
	case "*":
		return &object.FloatObject{Value: leftVal * rightVal}
	case "/":
		return &object.FloatObject{Value: leftVal / rightVal}
	case "<":
		return &object.BooleanObject{Value: leftVal < rightVal}
	case "<=":
		return &object.BooleanObject{Value: leftVal <= rightVal}
	case "==":
		return &object.BooleanObject{Value: leftVal == rightVal}
	case ">":
		return &object.BooleanObject{Value: leftVal > rightVal}
	case ">=":
		return &object.BooleanObject{Value: leftVal >= rightVal}
	default:
		return &object.ErrorObject{Message: fmt.Sprintf("unsupported operator '%s' for types %s, %s", op, left.Type(), right.Type())}
	}
}

func evalFloatInfixExpression(op string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.FloatObject).Value
	rightVal := right.(*object.FloatObject).Value

	switch op {
	case "+":
		return &object.FloatObject{Value: leftVal + rightVal}
	case "-":
		return &object.FloatObject{Value: leftVal - rightVal}
	case "*":
		return &object.FloatObject{Value: leftVal * rightVal}
	case "/":
		return &object.FloatObject{Value: leftVal / rightVal}
	case "<":
		return &object.BooleanObject{Value: leftVal < rightVal}
	case "<=":
		return &object.BooleanObject{Value: leftVal <= rightVal}
	case "==":
		return &object.BooleanObject{Value: leftVal == rightVal}
	case ">":
		return &object.BooleanObject{Value: leftVal > rightVal}
	case ">=":
		return &object.BooleanObject{Value: leftVal >= rightVal}
	default:
		return &object.ErrorObject{Message: fmt.Sprintf("unsupported operator '%s' for types %s, %s", op, left.Type(), right.Type())}
	}
}

func evalFunctionCall(functionCall *ast.FunctionCall, env *object.Environment) object.Object {
	if _, ok := env.Get(functionCall.Name); ok {
		return evalUserDefinedFun(functionCall, env)
	} else if _, ok := stdlib.BuiltInFuns[functionCall.Name]; ok {
		return evalBuiltInFun(functionCall, env)
	}

	fmt.Printf("function %s is not defined\n", functionCall.Name)
	os.Exit(1)

	return nil
}

func evalUserDefinedFun(functionCall *ast.FunctionCall, env *object.Environment) object.Object {
	obj, ok := env.Get(functionCall.Name)
	if !ok {
		fmt.Printf("function %s is not defined\n", functionCall.Name)
		os.Exit(1)
	}

	function := obj.(*object.FunctionObject)

	if len(functionCall.Args) != len(function.Args) {
		fmt.Printf("expected %d arguments for function %s, received %d\n", len(function.Args), functionCall.Name, len(functionCall.Args))
		os.Exit(1)
	}

	childEnv := object.CreateChildEnvironment(env)

	// assign function args as values in child environment
	for i := range function.Args {
		childEnv.Set(function.Args[i], Eval(functionCall.Args[i], env), true)
	}

	// evaluate each statement in the function
	for i := range function.Statements {
		res := Eval(function.Statements[i], childEnv)

		switch function.Statements[i].(type) {
		case *ast.ReturnStatement:
			return res
		}
	}

	// print out the state of the program
	// fmt.Printf("child env values for function call %s\n", functionCall)
	// for k, v := range childEnv.GetEnvMap() {
	// 	fmt.Printf("%s: %s\n", k, v.ToString())
	// }
	// fmt.Println("------------------------------")

	return &object.NullObject{}
}

func evalBuiltInFun(functionCall *ast.FunctionCall, env *object.Environment) object.Object {
	args := []object.Object{}

	for i := range functionCall.Args {
		args = append(args, Eval(functionCall.Args[i], env))
	}

	return stdlib.BuiltInFuns[functionCall.Name](args...)
}

func evalObjFunCall(objFunCall *ast.ObjectFunctionExpression, env *object.Environment) object.Object {
	args := []object.Object{Eval(objFunCall.Object, env)}
	fnCall := objFunCall.Function.(*ast.FunctionCall)

	for i := range fnCall.Args {
		args = append(args, Eval(fnCall.Args[i], env))
	}

	if fn, ok := stdlib.BuiltInFuns[fnCall.Name]; ok {
		return fn(args...)
	}

	return &object.ErrorObject{Message: fmt.Sprintf("function %s is not defined\n", fnCall.Name)}
}

func evalArrayExpression(items []ast.Expression, env *object.Environment) object.Object {
	arr := &object.ArrayObject{Items: []object.Object{}}

	for i := range items {
		arr.Items = append(arr.Items, Eval(items[i], env))
	}

	return arr
}

func evalArrayIndexExpression(arrIdxExpr *ast.ArrayIndexExpression, env *object.Environment) object.Object {
	arr, ok := Eval(arrIdxExpr.Arr, env).(*object.ArrayObject)

	if !ok {
		return &object.ErrorObject{Message: fmt.Sprintf("cannot index object of type %s", arr.Type())}
	}

	idx, ok := Eval(arrIdxExpr.Index, env).(*object.IntegerObject)

	if !ok {
		return &object.ErrorObject{Message: fmt.Sprintf("cannot use object of type %s as index", idx.Type())}
	}

	if int(idx.Value) >= len(arr.Items) {
		return &object.ErrorObject{Message: "array index out of bounds"}
	}

	return arr.Items[idx.Value]
}
