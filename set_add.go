package mpq

// Add adds the MPQ [Archive] specified at path to the [Set].
func (set *Set) Add(path string) (err error) {
	archive, err := Open(path)
	if err != nil {
		return
	}

	set.archives = append(set.archives, archive)
	return
}
