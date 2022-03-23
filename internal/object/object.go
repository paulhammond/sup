package object

import (
	"io/fs"
	"sort"
)

type Set map[string]Object

func (s Set) Paths() []string {
	keys := make([]string, len(s))

	i := 0
	for k := range s {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	return keys
}

type Object interface {
}

type File struct {
	path string
	fs   *fs.FS
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
			fs:   &filesystem,
		}

		return nil
	})

	return &set, err
}
