package object_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/paulhammond/sup/internal/object"
)

func TestSetPaths(t *testing.T) {
	set, err := object.FS(os.DirFS("../../testdata"))
	ok(t, err, "New")

	actual := set.Paths()
	expected := []string{"a.txt", "b.txt", "sub/a.txt"}
	if !(reflect.DeepEqual(actual, expected)) {
		t.Errorf("Wrong Paths:\ngot %#v\nexp %#v", actual, expected)
	}

}

func ok(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: unexpected error: %s", msg, err.Error())
	}
}
