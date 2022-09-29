package object

import "fmt"

type ArrayObject struct {
	Items []Object
}

func (a *ArrayObject) Type() string {
	return ARRAY_OBJ
}

func (a *ArrayObject) ToString() string {
	arrStr := ""

	for i := range a.Items {
		arrStr += a.Items[i].ToString()
		if i < len(a.Items)-1 {
			arrStr += ","
		}
	}

	return fmt.Sprintf("[%s]", arrStr)
}
