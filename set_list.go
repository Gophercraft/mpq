package mpq

import (
	"sort"
	"strings"
)

type set_list struct {
	combined_list []string
	index         int
}

func (set *Set) combine_list() (err error) {
	// merge lists into a dictionary to avoid repetition
	var dict = map[string]string{}
	var list List
	for _, archive := range set.archives {
		list, err = archive.List()
		if err != nil {
			return err
		}
		for list.Next() {
			path := list.Path()
			dict[strings.ToLower(path)] = path
		}
		list.Close()
	}
	// combine dictionary into a sorted slice of strings
	set.combined_list = make([]string, len(dict))
	i := 0
	for _, v := range dict {
		set.combined_list[i] = v
		i++
	}
	dict = nil
	sort.Strings(set.combined_list)
	return
}

func (list *set_list) Next() bool {
	return list.index < len(list.combined_list)
}

func (list *set_list) Path() string {
	path := list.combined_list[list.index]
	list.index++
	return path
}

func (list *set_list) Close() error {
	return nil
}

// (Slow)
// Returns a combined List for all the archives in the set.
// Call this after you've loaded all your Archives,
// the results are then frozen in memory and fast after the first slow call
func (set *Set) List() (list List, err error) {
	if len(set.combined_list) == 0 {
		if err = set.combine_list(); err != nil {
			return
		}
	}

	list = &set_list{
		index:         0,
		combined_list: set.combined_list,
	}
	return
}
