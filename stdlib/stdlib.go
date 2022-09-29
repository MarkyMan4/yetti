package stdlib

import (
	"fmt"

	"github.com/MarkyMan4/yetti/object"
)

type BuiltIn func(args ...object.Object) object.Object

var BuiltInFuns = map[string]BuiltIn{
	"print":  PrintFun,
	"substr": SubstringFun,
	"length": LengthFun,
}

func PrintFun(args ...object.Object) object.Object {
	// print each argument separated by space and ending with a newline
	for i := range args {
		fmt.Print(args[i].ToString())

		if i == len(args)-1 {
			fmt.Println()
		} else {
			fmt.Print(" ")
		}
	}

	return &object.NullObject{}
}

// TODO: handle indices out of bounds or invalid indices (e.g. index 2 < index 1)
func SubstringFun(args ...object.Object) object.Object {
	if args[0].Type() != object.STRING_OBJ {
		return &object.ErrorObject{Message: fmt.Sprintf("object of type %s has no function substring", args[0].Type())}
	}

	if len(args) < 2 || len(args) > 3 {
		return &object.ErrorObject{Message: "must provide one or two arguments to substring function"}
	}

	if len(args) >= 2 && args[1].Type() != object.INTEGER_OBJ {
		return &object.ErrorObject{Message: "arguments must be integers"}
	}

	if len(args) == 2 && args[2].Type() != object.INTEGER_OBJ {
		return &object.ErrorObject{Message: "arguments must be integers"}
	}

	strLit := args[0].(*object.StringObject).Value
	startIdx := args[1].(*object.IntegerObject).Value

	if len(args) == 2 {
		return &object.StringObject{Value: strLit[startIdx:]}
	}

	endIdx := args[2].(*object.IntegerObject).Value

	return &object.StringObject{Value: strLit[startIdx:endIdx]}
}

func LengthFun(args ...object.Object) object.Object {
	if args[0].Type() != object.STRING_OBJ && args[0].Type() != object.ARRAY_OBJ {
		return &object.ErrorObject{Message: fmt.Sprintf("object of type %s has no function length", args[0].Type())}
	}

	if len(args) > 1 {
		return &object.ErrorObject{Message: "length function takes no arguments"}
	}

	if args[0].Type() == object.STRING_OBJ {
		strObj := args[0].(*object.StringObject)
		return &object.IntegerObject{Value: int64(len(strObj.Value))}
	} else {
		arrObj := args[0].(*object.ArrayObject)
		return &object.IntegerObject{Value: int64(len(arrObj.Items))}
	}
}
