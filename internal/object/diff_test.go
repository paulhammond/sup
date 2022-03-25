package object_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/paulhammond/sup/internal/object"
)

type testObject string

var HashErr = errors.New("hash error")

func (o testObject) Hash() (string, error) {
	if o == "" {
		return "", HashErr
	}
	return string(o), nil
}

func (o testObject) Metadata() (*object.Metadata, error) {
	return nil, errors.New("not implemented")
}

func TestDiffBasic(t *testing.T) {

	old := object.Set{
		"same":  testObject("hash"),
		"diff":  testObject("old"),
		"extra": testObject("hash"),
		// won't check the hash on an extra object
		"extraerr": testObject(""),
	}
	new := object.Set{
		"same":  testObject("hash"),
		"diff":  testObject("new"),
		"added": testObject("hash"),
		// won't check the hash on an extra object
		"addederr": testObject(""),
	}

	toUpload, toDelete, err := old.Diff(new)
	ok(t, err, "Diff")

	expectedUploads := []string{"added", "addederr", "diff"}
	if !reflect.DeepEqual(expectedUploads, toUpload) {
		t.Fatalf("uploads wrong\ngot: %#v\nexp: %#v", toUpload, expectedUploads)
	}

	expectedDeletes := []string{"extra", "extraerr"}
	if !reflect.DeepEqual(expectedDeletes, toDelete) {
		t.Fatalf("deletes wrong\ngot: %#v\nexp: %#v", toDelete, expectedDeletes)
	}

}

func TestDiffMatch(t *testing.T) {

	set := object.Set{
		"same":  testObject("hash"),
		"diff":  testObject("old"),
		"extra": testObject("hash"),
	}

	toUpload, toDelete, err := set.Diff(set)
	ok(t, err, "Diff")

	if exp := []string{}; !reflect.DeepEqual(exp, toUpload) {
		t.Errorf("uploads wrong\ngot: %#v\nexp: %#v", toUpload, exp)
	}
	if exp := []string{}; !reflect.DeepEqual(exp, toDelete) {
		t.Errorf("deletes wrong\ngot: %#v\nexp: %#v", toDelete, exp)
	}

}

func TestDiffError(t *testing.T) {

	good := object.Set{
		"same": testObject("hash"),
		"bad":  testObject("ok"),
	}
	bad := object.Set{
		"same": testObject("hash"),
		"bad":  testObject(""),
	}

	toUpload, toDelete, err := good.Diff(bad)
	if err != HashErr {
		t.Errorf("err wrong\ngot: %$v\nexp: %#v", err, HashErr)
	}
	if exp := []string{}; !reflect.DeepEqual(exp, toUpload) {
		t.Errorf("uploads wrong\ngot: %#v\nexp: %#v", toUpload, exp)
	}
	if exp := []string{}; !reflect.DeepEqual(exp, toDelete) {
		t.Errorf("deletes wrong\ngot: %#v\nexp: %#v", toDelete, exp)
	}

	toUpload, toDelete, err = bad.Diff(good)
	if err != HashErr {
		t.Errorf("err wrong\ngot: %$v\nexp: %#v", err, HashErr)
	}

	if exp := []string{}; !reflect.DeepEqual(exp, toUpload) {
		t.Errorf("uploads wrong\ngot: %#v\nexp: %#v", toUpload, exp)
	}
	if exp := []string{}; !reflect.DeepEqual(exp, toDelete) {
		t.Errorf("deletes wrong\ngot: %#v\nexp: %#v", toDelete, exp)
	}

}
