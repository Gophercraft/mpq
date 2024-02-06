package mpq

import (
	"fmt"

	"github.com/Gophercraft/mpq/crypto"
	"github.com/Gophercraft/mpq/info"
)

func (archive *Archive) hash_lookup_name(name string) (index int, entry *info.HashTableEntry, err error) {
	// Convert name string into hashes
	name_1 := crypto.HashString(name, info.HashNameA)
	name_2 := crypto.HashString(name, info.HashNameB)

	// Lookup hashes in table
	index, entry, err = archive.search_hash_table(name_1, name_2)
	if err != nil {
		err = fmt.Errorf("mpq: could not find file named '%s' in archive '%s' (%s)", name, archive.path, err)
		return
	}

	return
}
