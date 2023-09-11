package pkg

import (
	"os"
)

// type FileSystem interface {
// 	Open(name string) (File, error)
// 	Stat(name string) (os.FileInfo, error)
// }
//
// type File interface {
// 	Stat() (os.FileInfo, error)
// }

type FS struct{}

func (*FS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (*FS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
