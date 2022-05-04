package object_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/paulhammond/sup/internal/object"
)

func TestSetPaths(t *testing.T) {
	dir := t.TempDir()

	err := os.MkdirAll(dir+"/sub", 0755)
	ok(t, err, "Mkdir sub")

	err = os.MkdirAll(dir+"/.git", 0755)
	ok(t, err, "Mkdir .git")

	err = os.WriteFile(dir+"/a.txt", []byte("a"), 0644)
	ok(t, err, "WriteFile a.txt")

	err = os.WriteFile(dir+"/b.txt", []byte("b"), 0644)
	ok(t, err, "WriteFile b.txt")

	err = os.WriteFile(dir+"/sub/a.txt", []byte("sub/a"), 0644)
	ok(t, err, "WriteFile a.txt")

	err = os.WriteFile(dir+"/.git/a.txt", []byte("sub/a"), 0644)
	ok(t, err, "WriteFile /git/a.txt")

	set, err := object.FS(os.DirFS(dir), false)
	ok(t, err, "New")

	actual := set.Paths()
	expected := []string{".git/a.txt", "a.txt", "b.txt", "sub/a.txt"}
	if !(reflect.DeepEqual(actual, expected)) {
		t.Errorf("Wrong Paths:\ngot %#v\nexp %#v", actual, expected)
	}

	set, err = object.FS(os.DirFS(dir), true)
	ok(t, err, "New")

	actual = set.Paths()
	expected = []string{"a.txt", "b.txt", "sub/a.txt"}
	if !(reflect.DeepEqual(actual, expected)) {
		t.Errorf("Wrong Paths:\ngot %#v\nexp %#v", actual, expected)
	}

}

func TestFileHash(t *testing.T) {
	set, err := object.FS(os.DirFS("testdata"), true)
	ok(t, err, "New")

	obj := set["a.txt"]

	hash, err := obj.Hash()
	ok(t, err, "Hash")

	exp := object.Hash{
		Size: 1,
		Hash: "0cc175b9c0f1b6a831c399e269772661",
	}
	if hash != nil && !reflect.DeepEqual(*hash, exp) {
		t.Errorf("Wrong Hash:\ngot %v\nexp %v", hash, exp)
	}

}

func ok(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: unexpected error: %s", msg, err.Error())
	}
}
