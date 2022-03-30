package filter

import "github.com/paulhammond/sup/internal/object"

func Filter(set *object.Set, debug DebugFunc) error {

	for _, p := range (*set).Paths() {
		err := detectType(p, (*set)[p], debug)
		if err != nil {
			return err
		}
	}

	return nil
}

type DebugFunc func(string, ...any)
