package filter

import (
	"os"
	"testing"

	"github.com/paulhammond/sup/internal/object"
)

func TestDetectType(t *testing.T) {

	objects, err := object.FS(os.DirFS("testdata/detect_type"))
	ok(t, err, "New")

	expected := map[string]string{
		"extension.txt":    "text/plain; charset=utf-8",
		"extension.html":   "text/html; charset=utf-8",
		"contents-html":    "text/html; charset=utf-8",
		"contents-unknown": "application/octet-stream",
	}
	for path, exp := range expected {
		err = detectType(path, objects[path])
		ok(t, err, "detectType "+path)

		m, err := objects[path].Metadata()
		ok(t, err, "metdata "+path)

		if m.ContentType == nil {
			t.Errorf("detectType %s:\ngot %#v\nexp %q", path, m.ContentType, exp)
		} else if *m.ContentType != exp {
			t.Errorf("detectType %s:\ngot %q\nexp %q", path, *m.ContentType, exp)
		}
	}
}

func ok(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: unexpected error: %s", msg, err.Error())
	}
}
