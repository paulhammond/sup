package filter

import (
	"strings"

	"github.com/paulhammond/sup/internal/object"
)

func ignoreDotfiles(set object.Set, debug DebugFunc) error {
	for _, path := range set.Paths() {
		ignore := isDotFile(path)

		if ignore {
			debug("dotfile [%s] ignoring dotfile", path)
			delete(set, path)
			continue
		}
	}
	return nil
}

func isDotFile(path string) bool {
	split := strings.Split(path, "/")
	for i, part := range split {
		if strings.HasPrefix(part, ".") {
			if i == 0 && part == ".well-known" {
				continue
			}
			return true
		}
	}
	return false
}
