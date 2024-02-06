package mpq

func (set *Set) Close() (err error) {
	for _, archive := range set.archives {
		if err = archive.Close(); err != nil {
			return
		}
	}

	return
}
