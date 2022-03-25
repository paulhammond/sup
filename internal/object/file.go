package object

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
)

type File struct {
	path     string
	fs       fs.FS
	metadata *Metadata
}

func (f File) Hash() (string, error) {
	var s string
	err := f.Open(func(r io.Reader) error {
		h := md5.New()
		size, err := io.Copy(h, r)
		if err != nil {
			return err
		}
		s = fmt.Sprintf("%d%x", size, h.Sum(nil))
		return nil
	})
	return s, err
}

func (f File) Open(fnc func(io.Reader) error) error {
	fd, err := f.fs.Open(f.path)
	if err != nil {
		return err
	}
	defer func() {
		e := fd.Close()
		if err == nil {
			err = e
		}
	}()

	return fnc(fd)
}

func (f File) Metadata() (*Metadata, error) {
	return f.metadata, nil
}

func FS(filesystem fs.FS) (Set, error) {
	set := Set{}

	err := fs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		set[path] = &File{
			path:     path,
			fs:       filesystem,
			metadata: &Metadata{},
		}

		return nil
	})

	return set, err
}
