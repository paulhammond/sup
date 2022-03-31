package filter

import (
	"strings"

	"github.com/paulhammond/sup/internal/cfg"
	"github.com/paulhammond/sup/internal/object"
)

func trim(config cfg.Config, set object.Set, debug DebugFunc) error {
	if len(config.TrimSuffix) == 0 {
		return nil
	}

	for _, path := range set.Paths() {
		orig := path
		for _, suffix := range config.TrimSuffix {
			path = strings.TrimSuffix(path, suffix)
		}
		if orig != path {
			debug("trim [%s] moved from %s", path, orig)
			set[path] = set[orig]
			delete(set, orig)
		}
	}

	return nil
}
