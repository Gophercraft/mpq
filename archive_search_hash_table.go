package mpq

import (
	"fmt"

	"github.com/Gophercraft/mpq/info"
)

func (archive *Archive) search_hash_table(name_1 uint32, name_2 uint32) (index int, result *info.HashTableEntry, err error) {
	for i := range archive.hash_table {
		hash_table_entry := &archive.hash_table[i]
		if hash_table_entry.Name1 == name_1 {
			if hash_table_entry.Name2 == name_2 {
				// Hash matches
				switch hash_table_entry.BlockIndex {
				case info.HashRemoved:
					continue
				case info.HashTerminator:
					break
				default:
					index = i
					result = hash_table_entry
					return
				}
			}
		}
	}

	err = fmt.Errorf("(0x%8x 0x%8x)", name_1, name_2)
	return
}
