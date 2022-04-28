package filter

import (
	"io"
	"testing"

	"github.com/paulhammond/sup/internal/cfg"
	"github.com/paulhammond/sup/internal/object"
)

func TestTrim(t *testing.T) {

	config := cfg.Config{TrimSuffix: []string{".trim"}}

	objects := object.Set{
		"one.trim": object.NewString("one"),
		"two.txt":  object.NewString("two"),
	}

	debug := newMockDebug()

	err := trim(config, objects, debug.debugFunc)
	ok(t, err, "trim")

	// was the trim file moved?
	if _, found := objects["one.trim"]; found {
		t.Fatalf("trim file not moved")
	}

	// did a new file get created?
	if _, found := objects["one"]; !found {
		t.Fatalf("trim file not moved")
	}

	// does the new file have the right contents?
	err = objects["one"].Open(func(r io.Reader) error {
		b, err := io.ReadAll(r)
		ok(t, err, "ReadAll")
		if got := string(b); got != "one" {
			t.Fatalf("trim file contents incorrect, got %q exp %q", got, "one")
		}
		return nil
	})
	ok(t, err, "open")

	// was the other file moved?
	if _, found := objects["two.txt"]; !found {
		t.Fatalf("non-trim file not left alone")
	}

	// did we log the right debugging info?
	expectedDebug := `trim [one] moved from one.trim` + "\n"
	if got := debug.String(); got != expectedDebug {
		t.Errorf("detectType debug:\ngot\n%s\nexp\n%s", got, expectedDebug)
	}

}
