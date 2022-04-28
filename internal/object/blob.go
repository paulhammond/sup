package object

import (
	"bytes"
	"io"
)

var _ Object = blob{}

type blob struct {
	contents []byte
	metadata *Metadata
}

func NewString(contents string) Object {
	return blob{contents: []byte(contents), metadata: &Metadata{}}
}

func Empty(metadata Metadata) Object {
	return blob{contents: []byte{}, metadata: &metadata}
}

func (o blob) Hash() (*Hash, error) {
	return GenerateHash(o)
}

func (o blob) Open(fnc func(io.Reader) error) error {
	r := bytes.NewReader(o.contents)
	return fnc(r)
}

func (o blob) Metadata() (*Metadata, error) {
	return o.metadata, nil
}
