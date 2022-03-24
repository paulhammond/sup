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

func TestFileHash(t *testing.T) {
	set, err := object.FS(os.DirFS("../../testdata"))
	ok(t, err, "New")

	obj := set["a.txt"]

	hash, err := obj.Hash()
	ok(t, err, "Hash")

	if exp := "10cc175b9c0f1b6a831c399e269772661"; hash != exp {
		t.Errorf("Wrong MD5:\ngot %q\nexp %q", hash, exp)
	}

}

func ok(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: unexpected error: %s", msg, err.Error())
	}
}
