package mpq

import (
	"fmt"

	"github.com/Gophercraft/mpq/info"
)

// The number of entries in the block table
func (archive *Archive) BlockTableCount() (count uint32) {
	count = archive.header.BlockTableSize
	return
}

// Return the block table entry at an index
func (archive *Archive) BlockTableIndex(block_table_index uint32) (block_table_entry *info.BlockTableEntry, err error) {
	if block_table_index >= uint32(len(archive.block_table)) {
		err = fmt.Errorf("mpq: block table entry index is out of bounds %d/%d", block_table_index, len(archive.block_table))
		return
	}

	block_table_entry = &archive.block_table[block_table_index]
	return
}
