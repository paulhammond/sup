package object

func (s Set) Diff(new Set) (Set, Set, error) {

	old := s // go likes all method recievers to have the same name
	// alternatively you could think of old as remote and new as local

	toDelete := make(Set, len(old))
	toUpload := make(Set, len(new))

	// start by checking all existing files to see if they should be deleted
	for k, o := range old {
		if _, ok := new[k]; !ok {
			toDelete[k] = o
		}
	}

	// now check all new files to see if they should be uploaded
	for k, f2 := range new {
		f1, ok := old[k]
		if !ok {
			toUpload[k] = f2
			continue
		}

		h1, err := f1.Hash()
		if err != nil {
			return Set{}, Set{}, err
		}
		h2, err := f2.Hash()
		if err != nil {
			return Set{}, Set{}, err
		}

		if h1 != h2 {
			toUpload[k] = f2
		}
	}

	return toUpload, toDelete, nil
}
