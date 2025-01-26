package mpq

import (
	"path/filepath"
)

// Opens a [Set] of archives using a list of glob patterns
func GlobSet(patterns ...string) (set *Set, err error) {
	set = NewSet()

	// glob archive paths
	var matched = map[string]bool{}
	var archive_list []string
	var matched_by_pattern []string
	var archive *Archive
	for _, pattern := range patterns {
		matched_by_pattern, err = filepath.Glob(pattern)
		if err != nil {
			return
		}

		archive_list = append(archive_list, matched_by_pattern...)
	}

	// attempt to open globbed archive paths
	for _, archive_path := range archive_list {
		if !matched[archive_path] {
			archive, err = Open(archive_path)
			if err != nil {
				return
			}

			set.archives = append(set.archives, archive)
			matched[archive_path] = true
		}
	}

	return
}
