package object

type StringObject struct {
	Value string
}

func (i *StringObject) Type() string {
	return STRING_OBJ
}

func (i *StringObject) ToString() string {
	return i.Value
}
