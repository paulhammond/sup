package object

import (
	"bytes"
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

func (o Blob) Hash() (*Hash, error) {
	return GenerateHash(o)
}

func (o Blob) Open(fnc func(io.Reader) error) error {
	r := bytes.NewReader(o.contents)
	return fnc(r)
}

func (o Blob) Metadata() (*Metadata, error) {
	return o.metadata, nil
}
