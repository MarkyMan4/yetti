package object

import "fmt"

type FileObject struct {
	FileName string
}

func (f *FileObject) Type() string {
	return FILE_OBJ
}

func (f *FileObject) ToString() string {
	return fmt.Sprintf("file object %s", f.FileName)
}
