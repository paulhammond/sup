package filter

import (
	"reflect"
	"strings"
	"testing"

	"github.com/paulhammond/sup/internal/cfg"
	"github.com/paulhammond/sup/internal/object"
)

func TestDotfiles(t *testing.T) {

	config := cfg.Config{}

	blank := object.NewString("")
	objects := object.Set{
		"file":                       blank,
		"dir/file":                   blank,
		"dir/.hidden":                blank,
		".hidden":                    blank,
		".hidden/file":               blank,
		".well-known/dnt-policy.txt": blank,
		".well-known/.hidden":        blank,
		"not-well-known/.well-known/dnt-policy.txt": blank,
	}

	debug := newMockDebug()

	err := ignoreDotfiles(config, objects, debug.debugFunc)
	ok(t, err, "ignore")

	exp := []string{
		".well-known/dnt-policy.txt",
		"dir/file",
		"file",
	}
	if !reflect.DeepEqual(exp, objects.Paths()) {
		t.Errorf("paths wrong\ngot: %#v\nexp: %#v", objects.Paths(), exp)
	}

	expectedDebug := `
dotfile [.hidden] ignoring dotfile
dotfile [.hidden/file] ignoring dotfile
dotfile [.well-known/.hidden] ignoring dotfile
dotfile [dir/.hidden] ignoring dotfile
dotfile [not-well-known/.well-known/dnt-policy.txt] ignoring dotfile
`
	expectedDebug = strings.TrimPrefix(expectedDebug, "\n")
	if got := debug.String(); got != expectedDebug {
		t.Errorf("dotfile debug:\ngot\n%s\nexp\n%s", got, expectedDebug)
	}

}

func TestDotfilesIncluded(t *testing.T) {

	config := cfg.Config{IncludeDotfiles: true}

	blank := object.NewString("")
	objects := object.Set{
		"file":    blank,
		".hidden": blank,
	}

	debug := newMockDebug()

	err := ignoreDotfiles(config, objects, debug.debugFunc)
	ok(t, err, "ignore")

	exp := []string{
		".hidden",
		"file",
	}
	if !reflect.DeepEqual(exp, objects.Paths()) {
		t.Errorf("paths wrong\ngot: %#v\nexp: %#v", objects.Paths(), exp)
	}

	expectedDebug := ""
	expectedDebug = strings.TrimPrefix(expectedDebug, "\n")
	if got := debug.String(); got != expectedDebug {
		t.Errorf("dotfile debug:\ngot\n%s\nexp\n%s", got, expectedDebug)
	}

}
