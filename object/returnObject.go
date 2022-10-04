package object

type ReturnObject struct {
	Value Object
}

func (f *ReturnObject) Type() string {
	return RETURN_OBJ
}

func (f *ReturnObject) ToString() string {
	return f.Value.ToString()
}
