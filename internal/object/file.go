package object

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
)

type File struct {
	path string
	fs   fs.FS
}

func (f File) Hash() (string, error) {
	fd, err := f.fs.Open(f.path)
	if err != nil {
		return "", err
	}
	defer func() {
		e := fd.Close()
		if err == nil {
			err = e
		}
	}()

	h := md5.New()
	var size int64
	if size, err = io.Copy(h, fd); err != nil {
		return "", err
	}
	return fmt.Sprintf("%d%x", size, h.Sum(nil)), nil
}

func FS(filesystem fs.FS) (*Set, error) {
	set := Set{}

	err := fs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		set[path] = &File{
			path: path,
			fs:   filesystem,
		}

		return nil
	})

	return &set, err
}
