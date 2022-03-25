package filter

import "github.com/paulhammond/sup/internal/object"

func Filter(set *object.Set) error {

	for p, o := range *set {
		err := detectType(p, o)
		if err != nil {
			return err
		}
	}

	return nil
}
