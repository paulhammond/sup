package remote_test

import (
	"testing"

	"github.com/paulhammond/sup/internal/remote"
)

func TestFakeHash(t *testing.T) {
	dir := t.TempDir()
	err := remote.CreateFake(dir + "/tmp.db")
	ok(t, err, "CreateFake")
	r, err := remote.Open(dir + "/tmp.db")
	ok(t, err, "Open")
	defer r.Close()

	set, err := r.Set()
	ok(t, err, "Set")

	obj := set["a.txt"]
	hash, err := obj.Hash()
	ok(t, err, "Hash")

	if exp := "2a1d0c6e83f027327d8461063f4ac58a6"; hash != exp {
		t.Errorf("Wrong MD5:\ngot %q\nexp %q", hash, exp)
	}
}

func ok(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: unexpected error: %s", msg, err.Error())
	}
}
