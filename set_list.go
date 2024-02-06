package mpq

import (
	"sort"
)

type set_list struct {
	combined_list []string
	index         int
}

func (set *Set) combine_list() (err error) {
	// merge lists into a dictionary to avoid repetition
	type nub struct{}
	var n nub
	var dict = map[string]nub{}
	var list List
	for _, archive := range set.archives {
		list, err = archive.List()
		if err != nil {
			return err
		}
		for list.Next() {
			dict[list.Path()] = n
		}
		list.Close()
	}
	// combine dictionary into a sorted slice of strings
	set.combined_list = make([]string, len(dict))
	i := 0
	for k := range dict {
		set.combined_list[i] = k
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
