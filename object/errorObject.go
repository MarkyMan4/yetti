package object

type ErrorObject struct {
	Message string
}

func (i *ErrorObject) Type() string {
	return STRING_OBJ
}

func (i *ErrorObject) ToString() string {
	return i.Message
}
