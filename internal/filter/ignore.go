package filter

import (
	"github.com/bmatcuk/doublestar/v2"
	"github.com/paulhammond/sup/internal/cfg"
	"github.com/paulhammond/sup/internal/object"
)

func ignore(cfg cfg.Config, set object.Set, debug DebugFunc) error {
	for _, path := range set.Paths() {
		ignore, err := isIgnored(cfg.Ignore, path, debug)

		if err != nil {
			return err
		}
		if ignore {
			delete(set, path)
			continue
		}
	}
	return nil
}

func isIgnored(ignore []string, path string, debug DebugFunc) (bool, error) {

	for _, pattern := range ignore {
		match, err := doublestar.Match(pattern, path)
		if err != nil {
			return false, err
		}
		if match {
			debug("ignore [%s] matches %q ignoring", path, pattern)
			return true, nil
		}
	}
	return false, nil
}
