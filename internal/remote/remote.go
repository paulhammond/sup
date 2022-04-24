package remote

import (
	"context"
	"strings"
	"time"

	"github.com/paulhammond/sup/internal/object"
)

var Timer = func() func() time.Duration {
	t0 := time.Now()
	return func() time.Duration {
		return time.Now().Sub(t0)
	}
}

type Remote interface {
	Close() error
	Set(context.Context) (object.Set, error)
	Upload(context.Context, object.Set, func(Event)) error
}

func Open(ctx context.Context, spec string) (Remote, error) {
	if strings.HasPrefix(spec, "s3://") {
		return openS3(ctx, spec)
	}
	return openFake(spec)
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
