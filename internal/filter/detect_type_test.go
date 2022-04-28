package filter

import (
	"testing"

	"github.com/paulhammond/sup/internal/object"
)

func TestDetectType(t *testing.T) {

	objects := object.Set{
		"contents-html":    object.NewString("<html><p>this is a html file</p>"),
		"contents-unknown": object.NewString("foo"),
		"extension.html":   object.NewString(""),
		"extension.txt":    object.NewString(""),
	}

	expected := [][2]string{
		{"extension.txt", "text/plain; charset=utf-8"},
		{"extension.html", "text/html; charset=utf-8"},
		{"contents-html", "text/html; charset=utf-8"},
		{"contents-unknown", "application/octet-stream"},
	}

	debug := newMockDebug()

	err := detectType(objects, debug.debugFunc)
	ok(t, err, "detectType")

	for _, tt := range expected {
		path, exp := tt[0], tt[1]

		m, err := objects[path].Metadata()
		ok(t, err, "metadata "+path)

		if m.ContentType == nil {
			t.Errorf("detectType %s:\ngot %#v\nexp %q", path, m.ContentType, exp)
		} else if *m.ContentType != exp {
			t.Errorf("detectType %s:\ngot %q\nexp %q", path, *m.ContentType, exp)
		}
	}

	expectedDebug := `
detecttype [contents-html] detected "text/html; charset=utf-8" via contents
detecttype [contents-unknown] detected "application/octet-stream" via contents
detecttype [extension.html] detected "text/html; charset=utf-8" via extension
detecttype [extension.txt] detected "text/plain; charset=utf-8" via extension
`[1:] // trim leading newline

	if got := debug.String(); got != expectedDebug {
		t.Errorf("detectType debug:\ngot\n%q\nexp\n%q", got, expectedDebug)
	}
}

func ok(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: unexpected error: %s", msg, err.Error())
	}
}
