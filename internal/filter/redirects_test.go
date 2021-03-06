package filter

import (
	"strings"
	"testing"

	"github.com/paulhammond/sup/internal/cfg"
	"github.com/paulhammond/sup/internal/object"
)

func TestRedirectFiles(t *testing.T) {

	config := cfg.Config{Redirects: true}

	objects := object.Set{
		"hello.redirect":   object.NewString("https://www.example.com/"),
		"newline.redirect": object.NewString("https://www.example.com/\n"),
		"not_redirect.txt": object.NewString("hello"),
	}

	debug := newMockDebug()

	err := processRedirect(config, objects, debug.debugFunc)
	ok(t, err, "processRedirect")

	// was the first redirect file moved?
	if _, found := objects["hello.redirect"]; found {
		t.Fatalf("redirect file not moved")
	}

	// did a new file get created?
	if _, found := objects["hello"]; !found {
		t.Fatalf("redirect file not moved")
	}

	metadata, err := objects["hello"].Metadata()
	ok(t, err, "metadata")
	checkStringRef(t, metadata.WebsiteRedirectLocation, str("https://www.example.com/"), "redirect location")

	// was the second redirect file moved?
	if _, found := objects["newline.redirect"]; found {
		t.Fatalf("redirect file not moved")
	}

	// did a new file get created?
	if _, found := objects["newline"]; !found {
		t.Fatalf("redirect file not moved")
	}

	metadata, err = objects["newline"].Metadata()
	ok(t, err, "metadata")
	checkStringRef(t, metadata.WebsiteRedirectLocation, str("https://www.example.com/"), "redirect location")

	// was the other file moved?
	if _, found := objects["not_redirect.txt"]; !found {
		t.Fatalf("non-redirect file not left alone")
	}
	metadata, err = objects["not_redirect.txt"].Metadata()
	ok(t, err, "metadata")
	checkStringRef(t, metadata.WebsiteRedirectLocation, nil, "non-redirect location")

	// did we log the right debugging info?
	expectedDebug := `
redirect [hello] created redirect to "https://www.example.com/"
redirect [newline] created redirect to "https://www.example.com/"
`
	expectedDebug = strings.TrimPrefix(expectedDebug, "\n")
	if got := debug.String(); got != expectedDebug {
		t.Errorf("detectType debug:\ngot\n%s\nexp\n%s", got, expectedDebug)
	}

}

func TestRedirectFilesDisabled(t *testing.T) {

	config := cfg.Config{Redirects: false}

	objects := object.Set{
		"hello.redirect":   object.NewString("https://www.example.com/"),
		"not_redirect.txt": object.NewString("hello"),
	}

	debug := newMockDebug()

	err := processRedirect(config, objects, debug.debugFunc)
	ok(t, err, "processRedirect")

	// was the redirect file moved?
	if _, found := objects["hello.redirect"]; !found {
		t.Fatalf("redirect file moved")
	}
	metadata, err := objects["not_redirect.txt"].Metadata()
	ok(t, err, "metadata")
	checkStringRef(t, metadata.WebsiteRedirectLocation, nil, "non-redirect location")

	// was the other file moved?
	if _, found := objects["not_redirect.txt"]; !found {
		t.Fatalf("non-redirect file not left alone")
	}
	metadata, err = objects["not_redirect.txt"].Metadata()
	ok(t, err, "metadata")
	checkStringRef(t, metadata.WebsiteRedirectLocation, nil, "non-redirect location")

	// did we log the right debugging info?
	expectedDebug := ""
	if got := debug.String(); got != expectedDebug {
		t.Errorf("detectType debug:\ngot\n%s\nexp\n%s", got, expectedDebug)
	}

}
