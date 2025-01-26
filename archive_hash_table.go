package mpq

import (
	"fmt"

	"github.com/Gophercraft/mpq/info"
)

// The number of entries in the hash table
func (archive *Archive) HashTableCount() (count uint32) {
	count = archive.header.HashTableSize
	return
}

// Return the hash table entry at the index
func (archive *Archive) HashTableIndex(index uint32) (hash_entry *info.HashTableEntry, err error) {
	if index >= uint32(len(archive.hash_table)) {
		err = fmt.Errorf("mpq: hash entry index out of bounds")
		return
	}

	hash_entry = &archive.hash_table[index]
	return
}
