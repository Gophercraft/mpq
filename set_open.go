package mpq

import "fmt"

// Open attempts to open the highest-order file contained within the [Set]
func (set *Set) Open(path string) (file *File, err error) {
	for i := len(set.archives) - 1; i >= 0; i-- {
		archive := set.archives[i]
		file, err = archive.Open(path)
		if err == nil {
			return
		}
	}

	return nil, fmt.Errorf("mpq: File '%s' not found in Set of %d Archives", path, len(set.archives))
}
