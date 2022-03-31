package filter

import (
	"errors"
	"io"
	"testing"

	"github.com/paulhammond/sup/internal/cfg"
	"github.com/paulhammond/sup/internal/object"
)

func TestAddMetadata(t *testing.T) {
	cfg := cfg.Config{
		Metadata: []cfg.Metadata{
			{Pattern: "**/*.txt", ContentType: str("text/plain; charset=utf-8")},
		},
	}

	debug := newMockDebug()

	tests := map[string]struct {
		ContentType *string
	}{
		"foo.png": {},
		"foo.txt": {ContentType: str("text/plain; charset=utf-8")},
	}

	for path, exp := range tests {
		o := testObject{m: &object.Metadata{}}
		err := addMetadata(cfg, path, o, debug.debugFunc)
		ok(t, err, "addMetadata "+path)

		m, err := o.Metadata()
		ok(t, err, "Metadata "+path)
		checkStringRef(t, m.ContentType, exp.ContentType, path+" ContentType")

	}

	expectedDebug := `metadata [foo.txt] matches "**/*.txt" set ContentType "text/plain; charset=utf-8"` + "\n"
	if got := debug.String(); got != expectedDebug {
		t.Errorf("detectType debug:\ngot\n%s\nexp\n%s", got, expectedDebug)
	}

}

func str(v string) *string {
	return &v
}

func checkStringRef(t *testing.T, got *string, expected *string, msg string) {
	t.Helper()
	if expected == nil {
		if got != nil {
			t.Fatalf("%s: got %q, expected nil", msg, *got)
		}
	} else {
		if got == nil {
			t.Fatalf("%s: got nil, expected %q", msg, *expected)
		} else {
			if *got != *expected {
				t.Fatalf("%s:\ngot %q\nexp %q", msg, *got, *expected)
			}
		}
	}
}

type testObject struct {
	m *object.Metadata
}

func (o testObject) Hash() (string, error) {
	return "", errors.New("not implemented")
}

func (o testObject) Metadata() (*object.Metadata, error) {
	return o.m, nil
}

func (o testObject) Open(func(io.Reader) error) error {
	return errors.New("not implemented")
}
