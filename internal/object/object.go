package object

import (
	"fmt"
	"io"
	"sort"
	"strings"
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

func (s Set) String() string {
	objects := make([]string, len(s))
	for i, k := range s.Paths() {
		hash, err := s[k].Hash()
		var str string
		if err != nil {
			str = "error"
		} else {
			str = fmt.Sprintf("%d.%d.%.6s", hash.Size, hash.PartSize, hash.Hash)
		}
		objects[i] = fmt.Sprintf("%s: %T<%s>", k, s[k], str)
	}

	return fmt.Sprintf("{\n\t%s\n}", strings.Join(objects, ",\n\t"))
}

type Object interface {
	Hash() (*Hash, error)
	Metadata() (*Metadata, error)
	Open(func(io.Reader) error) error
}

type Metadata struct {
	CacheControl            *string
	ContentType             *string
	WebsiteRedirectLocation *string
}
