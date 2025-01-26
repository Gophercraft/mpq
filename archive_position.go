package mpq

// Beginning returns the position of the archive's beginning
func (archive *Archive) Position() (pos int64) {
	return archive.archive_pos
}
