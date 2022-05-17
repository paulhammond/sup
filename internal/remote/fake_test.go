package remote_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/paulhammond/sup/internal/object"
	"github.com/paulhammond/sup/internal/remote"
)

func TestFakeDisallowed(t *testing.T) {
	ctx := context.Background()
	r, err := remote.Open(ctx, "tmp.db", false)
	if r != nil {
		t.Errorf("Unexpected remote\ngot %v\nexp nil", r)
	}
	if err == nil || err.Error() != "Unknown remote: tmp.db" {
		t.Errorf("Unexpected error\ngot %s", err)
	}
}

func TestFake(t *testing.T) {
	ctx := context.Background()

	dir := t.TempDir()
	err := remote.CreateFake(dir + "/tmp.db")
	ok(t, err, "CreateFake")
	r, err := remote.Open(ctx, dir+"/tmp.db", true)
	ok(t, err, "Open")
	defer r.Close()

	set, err := r.Set(ctx)
	ok(t, err, "Set")

	obj := set["a.txt"]

	hash, err := obj.Hash()
	ok(t, err, "Hash")
	exp := object.Hash{
		Size: 2,
		Hash: "a1d0c6e83f027327d8461063f4ac58a6",
	}
	if hash != nil && !reflect.DeepEqual(*hash, exp) {
		t.Errorf("Wrong Hash:\ngot %v\nexp %v", hash, exp)
	}

	metadata, err := obj.Metadata()
	ok(t, err, "Metadata")
	if exp := (&object.Metadata{ContentType: str("text/plain")}); !reflect.DeepEqual(exp, metadata) {
		t.Errorf("metadata wrong\ngot: %#v\nexp: %#v", metadata, exp)
	}
}

func ok(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: unexpected error: %s", msg, err.Error())
	}
}

func str(v string) *string {
	return &v
}
