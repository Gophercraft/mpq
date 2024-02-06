package mpq

func (set *Set) Add(path string) (err error) {
	archive, err := Open(path)
	if err != nil {
		return
	}

	set.archives = append(set.archives, archive)
	return
}
