package mpq

import "fmt"

// Return whether the Archive contains a hi-block table
//
// A hi-block table contains high parts of block table offsets, allowing MPQ archives larger than 4GB
func (archive *Archive) ContainsHiBlockTable() bool {
	return archive.header.HiBlockTablePos64 != 0
}

// returns the absolute position of the Archive's hi-block table
func (archive *Archive) get_hi_block_table_position() (absolute_pos uint64, err error) {
	absolute_pos = uint64(archive.archive_pos) + archive.header.HiBlockTablePos64
	return
}

// Return the higher order bits of a block table index
func (archive *Archive) HiBlockTableIndex(block_table_index uint32) (hi_block_position uint16, err error) {
	if block_table_index >= uint32(len(archive.hi_block_table)) {
		err = fmt.Errorf("mpq: hi block table index out of bounds")
		return
	}

	hi_block_position = archive.hi_block_table[block_table_index]
	return
}
