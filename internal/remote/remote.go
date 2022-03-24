package remote

import "github.com/paulhammond/sup/internal/object"

type Remote interface {
	Close() error
	Set() (object.Set, error)
}

func Open(p string) (Remote, error) {
	return openFake(p)
}
