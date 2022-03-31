package object

import (
	"io"
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
	Hash() (string, error)
	Metadata() (*Metadata, error)
	Open(func(io.Reader) error) error
}

type Metadata struct {
	ContentType             *string
	WebsiteRedirectLocation *string
}
