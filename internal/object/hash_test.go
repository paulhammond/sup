package object

import (
	"bytes"
	"io"
	"testing"
)

type TestObject struct {
	Contents string
	hash     *Hash
	Err      error
}

func (o TestObject) Hash() (*Hash, error) {
	if o.hash != nil || o.Err != nil {
		return o.hash, o.Err
	}
	return GenerateHash(o)
}

func (o TestObject) Open(fnc func(io.Reader) error) error {
	return fnc(bytes.NewReader([]byte(o.Contents)))
}

func (o TestObject) Metadata() (*Metadata, error) {
	panic("unimplemented")
}

func TestHashMatch(t *testing.T) {

	tests := []struct {
		name     string
		exp      bool
		contents string
		hash     *Hash
	}{
		{"simple", true, "hellothere", nil},
		{"nomatch", false, "hi", nil},
		{"hash", true, "hellothere", &Hash{
			Size:     10,
			PartSize: 0,
			Hash:     "c6f7c372641dd25e0fddf0215375561f",
		}},
		{"multipart", true, "hellothere", &Hash{
			Size:     10,
			PartSize: 6,
			Hash:     "405bee0ce44e8f5bc2cc560cda057716-2",
		}},
		{"multipartnomatch", false, "hellothere", &Hash{
			Size:     10,
			PartSize: 6,
			Hash:     "405bee0ce44e8f5bc2cc560cda057717-2",
		}},
		{"mismatchsize", false, "hellothere", &Hash{
			Size:     11,
			PartSize: 0,
			Hash:     "c6f7c372641dd25e0fddf0215375561f",
		}},
	}

	o := TestObject{Contents: "hellothere"}
	for _, tt := range tests {
		testObject := TestObject{Contents: tt.contents, hash: tt.hash}

		got, err := matchHash(testObject, o)
		if err != nil {
			t.Errorf("%s: unexpected error: %s", tt.name, err.Error())
			continue
		}
		if got != tt.exp {
			t.Errorf("%s: unexpected result, got %t exp %t", tt.name, got, tt.exp)
		}
	}
}
