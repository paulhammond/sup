package filter

import (
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

	set := object.Set{
		"foo.png": object.NewBlob([]byte{}, object.Metadata{}),
		"foo.txt": object.NewBlob([]byte{}, object.Metadata{}),
	}

	tests := map[string]struct {
		ContentType *string
	}{
		"foo.png": {},
		"foo.txt": {ContentType: str("text/plain; charset=utf-8")},
	}

	err := addMetadata(cfg, set, debug.debugFunc)
	ok(t, err, "addMetadata")

	for path, exp := range tests {
		o := set[path]
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
