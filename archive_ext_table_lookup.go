package mpq

import (
	"fmt"
	"strings"

	"github.com/Gophercraft/mpq/info"
	"github.com/Gophercraft/mpq/jenkins"
)

func (archive *Archive) ext_table_lookup(name string) (bet_index uint32, err error) {
	// Calculate 64-bit hash of the file name
	name_hash_64 := jenkins.Hash64([]byte(strings.ToLower(name)))
	name_hash := info.HetTableLookupValue(archive.HetTableHeader(), name_hash_64)
	// Split the file name hash into two parts:
	// NameHash1: The highest 8 bits of the name hash
	// NameHash2: File name hash limited to hash size
	name_hash_1 := uint8(name_hash >> (uint64(archive.het_table.header.NameHashBitSize - 8)))

	start_index := uint32(name_hash % uint64(archive.HetTableHeader().TotalEntryCount))
	index := start_index

	for archive.het_table.name_hashes[index] != 0x00 {
		if archive.het_table.name_hashes[index] == name_hash_1 {
			var possible_match_bet_index uint32
			possible_match_bet_index, err = archive.HetTableIndexBetTableIndex(index)
			if err != nil {
				return
			}

			if possible_match_bet_index >= archive.het_table.header.TotalEntryCount {
				err = fmt.Errorf("mpq: possible match is not within bounds of HET table")
				return
			}

			var name_hash_2 uint64
			name_hash_2, err = archive.BetTableNameHash2Index(uint32(possible_match_bet_index))
			if err != nil {
				return
			}

			possible_hash := info.BetTableMergeHashValue(archive.BetTableHeader(), name_hash_1, name_hash_2)
			if possible_hash == name_hash {
				bet_index = uint32(possible_match_bet_index)
				return
			}
		}

		index = (index + 1) % archive.HetTableHeader().TotalEntryCount
		if index == start_index {
			break
		}
	}

	err = fmt.Errorf("mpq: file not found in HET table")
	return
}
