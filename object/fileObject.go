package object

type FileObject struct {
	FileName string
	Mode     string // "r" or "w" TODO: add "a" for append mode
}

func (f *FileObject) Type() string {
	return FILE_OBJ
}

func (f *FileObject) ToString() string {
	return "file object"
}
