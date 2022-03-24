package object

import "sort"

func (s Set) Diff(new Set) ([]string, []string, error) {

	old := s // go likes all method recievers to have the same name
	// alternatively you could think of old as remote and new as local

	toDelete := make([]string, 0, len(old))
	toUpload := make([]string, 0, len(new))

	// start by checking all existing files to see if they should be deleted
	for k := range old {
		if _, ok := new[k]; !ok {
			toDelete = append(toDelete, k)
		}
	}

	// now check all new files to see if they should be uploaded
	for k, f2 := range new {
		f1, ok := old[k]
		if !ok {
			toUpload = append(toUpload, k)
			continue
		}

		h1, err := f1.Hash()
		if err != nil {
			return []string{}, []string{}, err
		}
		h2, err := f2.Hash()
		if err != nil {
			return []string{}, []string{}, err
		}

		if h1 != h2 {
			toUpload = append(toUpload, k)
		}
	}

	sort.Strings(toUpload)
	sort.Strings(toDelete)

	return toUpload, toDelete, nil
}
