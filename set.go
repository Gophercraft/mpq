package mpq

// Set contains a list of loaded MoPaQ [Archive]s
type Set struct {
	archives      []*Archive
	combined_list []string
}

// Returns an empty [Set]
func NewSet() (set *Set) {
	set = new(Set)
	return
}
