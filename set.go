package mpq

// Set contains a list of loaded MPQ Archives
type Set struct {
	archives      []*Archive
	combined_list []string
}

// Returns an empty Set
func NewSet() (set *Set) {
	set = new(Set)
	return
}
