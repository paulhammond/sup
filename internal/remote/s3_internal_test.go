package remote

import "testing"

func TestS3prefix(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{in: "/", out: ""},
		{in: "/foo", out: "foo/"},
		{in: "/foo/", out: "foo/"},
		{in: "foo", out: "foo/"},
		{in: "foo/", out: "foo/"},
		{in: "/foo/bar/baz", out: "foo/bar/baz/"},
		{in: "/foo/bar/baz/", out: "foo/bar/baz/"},
		{in: "foo/bar/baz", out: "foo/bar/baz/"},
		{in: "foo/bar/baz/", out: "foo/bar/baz/"},
	}
	for _, tt := range tests {
		if got := normalizePrefix(tt.in); got != tt.out {
			t.Errorf("normalizePrefix(%q), got %q exp %q", tt.in, got, tt.out)
		}
	}
}
