package object

import "fmt"

type BooleanObject struct {
	Value bool
}

func (i *BooleanObject) Type() string {
	return BOOLEAN_OBJ
}

func (i *BooleanObject) ToString() string {
	return fmt.Sprint(i.Value)
}
