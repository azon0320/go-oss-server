package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileSystem struct {
	path string
}

func (f *FileSystem) GetRootPath() string {
	return f.path
}

func (f *FileSystem) Path(subPath string) string {
	return fmt.Sprintf("%s%s%s",
		f.GetRootPath(),
		f.FileSeparator(),
		strings.Trim(subPath, f.FileSeparator()))
}

func (f *FileSystem) FileSeparator() string {
	return string(filepath.Separator)
}

func (f *FileSystem) Exists(sub string) bool {
	var exists = true
	var fullPath = f.Path(sub)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		exists = false
	}
	return exists
}

func (f *FileSystem) Create(sub string) error {
	var fullpath = f.Path(sub)
	file, err := os.Create(fullpath)
	if err == nil {
		file.Close()
	}
	return err
}

func (f *FileSystem) Delete(sub string) error {
	return os.Remove(f.Path(sub))
}

func CreateFileSystem(path string, perm os.FileMode) *FileSystem {
	var fileSeparator = string(filepath.Separator)
	fs := &FileSystem{
		path: strings.TrimRight(path, fileSeparator),
	}
	os.MkdirAll(fs.GetRootPath(), perm)
	return fs
}
