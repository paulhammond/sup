package object

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
)

var _ Object = Blob{}

type Blob struct {
	contents []byte
	metadata *Metadata
}

func NewBlob(contents []byte, metadata Metadata) Object {
	return Blob{contents: contents, metadata: &metadata}
}

func (o Blob) Hash() (string, error) {
	h := md5.Sum(o.contents)
	return fmt.Sprintf("%d%x", len(o.contents), h[:]), nil
}

func (o Blob) Open(fnc func(io.Reader) error) error {
	r := bytes.NewReader(o.contents)
	return fnc(r)
}

func (o Blob) Metadata() (*Metadata, error) {
	return o.metadata, nil
}
