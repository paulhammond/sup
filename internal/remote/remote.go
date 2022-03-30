package remote

import (
	"time"

	"github.com/paulhammond/sup/internal/object"
)

type Remote interface {
	Close() error
	Set() (object.Set, error)
	Upload(object.Set, func(Event)) error
}

func Open(p string) (Remote, error) {
	return openFake(p)
}

type Operation = int

const (
	Upload Operation = iota
	Download
)

type Event struct {
	Type     Operation
	Path     string
	Object   object.Object
	Duration time.Duration
}
