package object

import "fmt"

type FloatObject struct {
	Value float64
}

func (i *FloatObject) Type() string {
	return FLOAT_OBJ
}

func (i *FloatObject) ToString() string {
	return fmt.Sprint(i.Value)
}
