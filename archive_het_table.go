package mpq

import (
	"fmt"

	"github.com/Gophercraft/mpq/info"
)

func (archive *Archive) contains_het_table() bool {
	return archive.header.BetTablePos64 != 0
}

// The header of the HET table
func (archive *Archive) HetTableHeader() (header *info.HetTableHeader) {
	return &archive.het_table.header
}

// The number of entries in the HET table
func (archive *Archive) HetTableCount() (count uint32) {
	count = archive.het_table.header.TotalEntryCount
	return
}

// Return the name hash 1 value at an index in the HET table
func (archive *Archive) HetTableNameHash1Index(het_table_index uint32) (name_hash_1 uint8, err error) {
	if het_table_index >= archive.HetTableCount() {
		err = fmt.Errorf("mpq: het table index out of bounds")
		return
	}

	name_hash_1 = archive.het_table.name_hashes[het_table_index]

	return
}

// Return the BET table index of a HET table entry index
func (archive *Archive) HetTableIndexBetTableIndex(het_table_index uint32) (bet_table_index uint32, err error) {
	if het_table_index >= archive.HetTableCount() {
		err = fmt.Errorf("mpq: het table index out of bounds")
		return
	}

	index_bit_length := uint64(archive.het_table.header.IndexSizeTotal)

	var bet_table_index_64 uint64
	bet_table_index_64, err = archive.het_table.bet_indices.Uint(index_bit_length*uint64(het_table_index), uint8(archive.het_table.header.IndexSize))
	if err != nil {
		return
	}

	bet_table_index = uint32(bet_table_index_64)
	return
}
