package filter

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/paulhammond/sup/internal/cfg"
	"github.com/paulhammond/sup/internal/object"
)

func TestIgnore(t *testing.T) {

	objects, err := object.FS(os.DirFS("testdata/ignore"))
	ok(t, err, "New")

	config := cfg.Config{Ignore: []string{"one.*"}}
	debug := newMockDebug()

	err = ignore(config, objects, debug.debugFunc)
	ok(t, err, "ignore")

	if exp := []string{"two.txt"}; !reflect.DeepEqual(exp, objects.Paths()) {
		t.Errorf("paths wrong\ngot: %#v\nexp: %#v", objects.Paths(), exp)
	}

	// test the last set of paths
	expectedDebug := `
ignore [one.txt] matches "one.*" ignoring
`[1:] // trim leading newline

	if got := debug.String(); got != expectedDebug {
		t.Errorf("detectType debug:\ngot\n%s\nexp\n%s", got, expectedDebug)
	}

}

func TestIsIgnored(t *testing.T) {

	tests := []struct {
		ignore []string
		paths  map[string]bool
	}{
		{
			ignore: []string{},
			paths: map[string]bool{
				"file": false,
			},
		},
		{
			ignore: []string{"foo"},
			paths: map[string]bool{
				"file":    false,
				"foo":     true,
				"foobar":  false,
				"foo/bar": false,
			},
		},
		{
			ignore: []string{"foo*"},
			paths: map[string]bool{
				"file":    false,
				"foo":     true,
				"foobar":  true,
				"foo/bar": false,
			},
		},
		{
			ignore: []string{"foo/*"},
			paths: map[string]bool{
				"file":        false,
				"foo":         false,
				"foobar":      false,
				"foo/bar":     true,
				"foo/bar/baz": false,
			},
		},
		{
			ignore: []string{"foo/**"},
			paths: map[string]bool{
				"file":        false,
				"foo":         false,
				"foobar":      false,
				"foo/bar":     true,
				"foo/bar/baz": true,
			},
		},
	}

	debug := newMockDebug()
	for _, tt := range tests {
		for path, exp := range tt.paths {
			name := fmt.Sprintf("%v %s", tt.ignore, path)
			got, err := isIgnored(tt.ignore, path, debug.debugFunc)
			ok(t, err, name)
			if got != exp {
				t.Errorf("%s isIgnored: got: %t exp: %t", name, got, exp)
			}
		}
	}

}
