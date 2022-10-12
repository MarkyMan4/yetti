package stdlib

import (
	"bufio"
	"fmt"
	"os"

	"github.com/MarkyMan4/yetti/object"
)

type BuiltIn func(args ...object.Object) object.Object

var BuiltInFuns = map[string]BuiltIn{
	"print":  PrintFun,
	"substr": SubstringFun,
	"length": LengthFun,
	"append": ArrayAppendFun,
	"string": StringFun,
	"input":  InputFun,
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

func InputFun(args ...object.Object) object.Object {
	if len(args) > 1 {
		return &object.ErrorObject{Message: fmt.Sprintf("input expects 0 or 1 arguments but received %d", len(args))}
	}

	if len(args) > 0 {
		if args[0].Type() != object.STRING_OBJ {
			return &object.ErrorObject{Message: fmt.Sprintf("input expects string argument but received object of type %s", args[0].Type())}
		}

		fmt.Print(args[0].ToString())
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return &object.StringObject{Value: scanner.Text()}
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

// get length of string or array
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

// append to an array
func ArrayAppendFun(args ...object.Object) object.Object {
	if args[0].Type() != object.ARRAY_OBJ {
		return &object.ErrorObject{Message: fmt.Sprintf("object of type %s has no function append", args[0].Type())}
	}

	if len(args) != 2 {
		return &object.ErrorObject{Message: "append takes exactly one argument"}
	}

	arr := args[0].(*object.ArrayObject)
	arr.Items = append(arr.Items, args[1])

	return arr
}

// convert object to string object
func StringFun(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.ErrorObject{Message: "string takes exactly one argument"}
	}

	return &object.StringObject{Value: args[0].ToString()}
}

func OpenFileFun(args ...object.Object) object.Object {
	if len(args) != 2 {
		return &object.ErrorObject{Message: "open takes exactly two arguments"}
	}

	if args[0].Type() != object.STRING_OBJ {
		return &object.ErrorObject{Message: "first argument must be a file name"}
	}

	if args[1].Type() != object.STRING_OBJ && args[1].ToString() != "r" && args[1].ToString() != "w" {
		return &object.ErrorObject{Message: "second argument must be \"r\" or \"w\""}
	}

	if _, err := os.Stat(args[0].ToString()); err != nil && args[1].ToString() != "w" {
		return &object.ErrorObject{Message: fmt.Sprintf("file %s does not exist", args[0].ToString())}
	}

	return &object.FileObject{
		FileName: args[0].ToString(),
		Mode:     args[1].ToString(),
	}
}
