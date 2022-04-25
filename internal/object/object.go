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
		if err != nil {
			hash = "error"
		}
		objects[i] = fmt.Sprintf("%s: %T<%.6s>", k, s[k], hash)
	}

	return fmt.Sprintf("{\n\t%s\n}", strings.Join(objects, ",\n\t"))
}

type Object interface {
	Hash() (string, error)
	Metadata() (*Metadata, error)
	Open(func(io.Reader) error) error
}

type Metadata struct {
	CacheControl            *string
	ContentType             *string
	WebsiteRedirectLocation *string
}
