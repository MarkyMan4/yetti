package object

const (
	INTEGER_OBJ  = "INTEGER"
	FLOAT_OBJ    = "FLOAT"
	STRING_OBJ   = "STRING"
	BOOLEAN_OBJ  = "BOOLEAN"
	ERROR_OBJ    = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
	ARRAY_OBJ    = "ARRAY"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN_OBJ"
	FILE_OBJ     = "FILE"
)

type Object interface {
	Type() string
	ToString() string
}
