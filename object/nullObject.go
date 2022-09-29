package object

type NullObject struct{}

func (i *NullObject) Type() string {
	return NULL_OBJ
}

func (i *NullObject) ToString() string {
	return "null"
}
