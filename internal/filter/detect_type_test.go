package filter

import (
	"os"
	"testing"

	"github.com/paulhammond/sup/internal/object"
)

func TestDetectType(t *testing.T) {

	objects, err := object.FS(os.DirFS("testdata/detect_type"))
	ok(t, err, "New")

	expected := [][2]string{
		{"extension.txt", "text/plain; charset=utf-8"},
		{"extension.html", "text/html; charset=utf-8"},
		{"contents-html", "text/html; charset=utf-8"},
		{"contents-unknown", "application/octet-stream"},
	}

	debug := newMockDebug()

	for _, tt := range expected {
		path, exp := tt[0], tt[1]

		err = detectType(path, objects[path], debug.debugFunc)
		ok(t, err, "detectType "+path)

		m, err := objects[path].Metadata()
		ok(t, err, "metdata "+path)

		if m.ContentType == nil {
			t.Errorf("detectType %s:\ngot %#v\nexp %q", path, m.ContentType, exp)
		} else if *m.ContentType != exp {
			t.Errorf("detectType %s:\ngot %q\nexp %q", path, *m.ContentType, exp)
		}
	}

	expectedDebug := `
detecttype [extension.txt] detected "text/plain; charset=utf-8" via extension
detecttype [extension.html] detected "text/html; charset=utf-8" via extension
detecttype [contents-html] detected "text/html; charset=utf-8" via contents
detecttype [contents-unknown] detected "application/octet-stream" via contents
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
