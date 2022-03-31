package filter

import (
	"github.com/paulhammond/sup/internal/cfg"
	"github.com/paulhammond/sup/internal/object"
)

func Filter(set *object.Set, cfg cfg.Config, debug DebugFunc) error {

	err := ignore(cfg, (*set), debug)
	if err != nil {
		return err
	}

	err = processRedirect(cfg, *set, debug)
	if err != nil {
		return err
	}

	for _, p := range (*set).Paths() {
		o := (*set)[p]
		err := addMetadata(cfg, p, o, debug)
		if err != nil {
			return err
		}

		err = detectType(p, o, debug)
		if err != nil {
			return err
		}
	}

	err = trim(cfg, *set, debug)
	if err != nil {
		return err
	}

	return nil
}

type DebugFunc func(string, ...any)
