package object_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/paulhammond/sup/internal/object"
)

var HashErr = errors.New("hash error")

func TestDiffBasic(t *testing.T) {

	old := object.Set{
		"same":  object.TestObject{Contents: "hash"},
		"diff":  object.TestObject{Contents: "old"},
		"extra": object.TestObject{Contents: "hash"},
		// won't check the hash on an extra object
		"extraerr": object.TestObject{Err: HashErr},
	}
	new := object.Set{
		"same":  object.TestObject{Contents: "hash"},
		"diff":  object.TestObject{Contents: "new"},
		"added": object.TestObject{Contents: "hash"},
		// won't check the hash on an extra object
		"addederr": object.TestObject{Err: HashErr},
	}

	toUpload, toDelete, err := old.Diff(new)
	ok(t, err, "Diff")

	expectedUploads := []string{"added", "addederr", "diff"}
	if !reflect.DeepEqual(expectedUploads, toUpload.Paths()) {
		t.Fatalf("uploads wrong\ngot: %#v\nexp: %#v", toUpload, expectedUploads)
	}

	expectedDeletes := []string{"extra", "extraerr"}
	if !reflect.DeepEqual(expectedDeletes, toDelete.Paths()) {
		t.Fatalf("deletes wrong\ngot: %#v\nexp: %#v", toDelete, expectedDeletes)
	}

}

func TestDiffMatch(t *testing.T) {

	set := object.Set{
		"same":  object.TestObject{Contents: "hash"},
		"diff":  object.TestObject{Contents: "old"},
		"extra": object.TestObject{Contents: "hash"},
	}

	toUpload, toDelete, err := set.Diff(set)
	ok(t, err, "Diff")

	if exp := []string{}; !reflect.DeepEqual(exp, toUpload.Paths()) {
		t.Errorf("uploads wrong\ngot: %#v\nexp: %#v", toUpload, exp)
	}
	if exp := []string{}; !reflect.DeepEqual(exp, toDelete.Paths()) {
		t.Errorf("deletes wrong\ngot: %#v\nexp: %#v", toDelete, exp)
	}

}

func TestDiffError(t *testing.T) {

	good := object.Set{
		"same": object.TestObject{Contents: "hash"},
		"bad":  object.TestObject{Contents: "ok"},
	}
	bad := object.Set{
		"same": object.TestObject{Contents: "hash"},
		"bad":  object.TestObject{Err: HashErr},
	}

	toUpload, toDelete, err := good.Diff(bad)
	if err != HashErr {
		t.Errorf("err wrong\ngot: %$v\nexp: %#v", err, HashErr)
	}
	if exp := []string{}; !reflect.DeepEqual(exp, toUpload.Paths()) {
		t.Errorf("uploads wrong\ngot: %#v\nexp: %#v", toUpload, exp)
	}
	if exp := []string{}; !reflect.DeepEqual(exp, toDelete.Paths()) {
		t.Errorf("deletes wrong\ngot: %#v\nexp: %#v", toDelete, exp)
	}

	toUpload, toDelete, err = bad.Diff(good)
	if err != HashErr {
		t.Errorf("err wrong\ngot: %$v\nexp: %#v", err, HashErr)
	}

	if exp := []string{}; !reflect.DeepEqual(exp, toUpload.Paths()) {
		t.Errorf("uploads wrong\ngot: %#v\nexp: %#v", toUpload, exp)
	}
	if exp := []string{}; !reflect.DeepEqual(exp, toDelete.Paths()) {
		t.Errorf("deletes wrong\ngot: %#v\nexp: %#v", toDelete, exp)
	}

}
