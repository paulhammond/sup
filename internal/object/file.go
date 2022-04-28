package object

import (
	"io"
	"io/fs"
)

var _ Object = File{}

type File struct {
	path     string
	fs       fs.FS
	metadata *Metadata
}

func (f File) Hash() (*Hash, error) {
	return GenerateHash(f)
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

func FS(filesystem fs.FS, skipGit bool) (Set, error) {
	set := Set{}

	err := fs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// make debug logs slightly clearer
		if skipGit && d.IsDir() && d.Name() == ".git" {
			return fs.SkipDir
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
