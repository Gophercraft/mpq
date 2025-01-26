package mpq

import (
	"fmt"

	"github.com/Gophercraft/mpq/info"
)

// The header of the BET table
func (archive *Archive) BetTableHeader() (bet_table_header *info.BetTableHeader) {
	bet_table_header = &archive.bet_table.header
	return
}

// The number of entries in the BET table
func (archive *Archive) BetTableCount() (bet_table_count uint32) {
	bet_table_count = archive.bet_table.header.EntryCount
	return
}

// Look up an entry in the BET table, copying its information into bet_table_entry
func (archive *Archive) BetTableEntryIndex(index uint32, bet_table_entry *info.BetTableEntry) (err error) {
	bet_table := &archive.bet_table

	entry_size := uint64(bet_table.header.TableEntrySize)

	entry_position := entry_size * uint64(index)

	if entry_position >= bet_table.entries.Len() || entry_position+entry_size > bet_table.entries.Len() {
		err = fmt.Errorf("mpq: BET table index is out of bounds")
		return
	}

	// Read the file position
	bet_table_entry.Position, err = bet_table.entries.Uint(entry_position+uint64(bet_table.header.BitIndex_FilePos), uint8(bet_table.header.BitCount_FilePos))
	if err != nil {
		return
	}

	bet_table_entry.FileSize, err = bet_table.entries.Uint(entry_position+uint64(bet_table.header.BitIndex_FileSize), uint8(bet_table.header.BitCount_FileSize))
	if err != nil {
		return
	}

	bet_table_entry.BlockSize, err = bet_table.entries.Uint(entry_position+uint64(bet_table.header.BitIndex_CompressedSize), uint8(bet_table.header.BitCount_CompressedSize))
	if err != nil {
		return
	}

	if bet_table.header.FlagCount != 0 {
		var flag_index uint64
		flag_index, err = bet_table.entries.Uint(entry_position+uint64(bet_table.header.BitIndex_FlagIndex), uint8(bet_table.header.BitCount_FlagIndex))
		if err != nil {
			return
		}
		if flag_index >= uint64(bet_table.header.FlagCount) {
			err = fmt.Errorf("mpq: bad flag count")
			return
		}
		bet_table_entry.FlagsIndex = uint32(flag_index)
	}

	return
}

// The number of unique file flags values in the BET table
func (archive *Archive) BetTableFileFlagsCount() (file_flags_count uint32) {
	file_flags_count = archive.BetTableHeader().FlagCount
	return
}

// Return file flags for an file flags index
func (archive *Archive) BetTableFileFlagsIndex(file_flags_index uint32) (file_flags info.FileFlag, err error) {
	if file_flags_index >= archive.BetTableFileFlagsCount() {
		err = fmt.Errorf("mpq: file flag index out of bounds")
		return
	}
	file_flags = archive.bet_table.file_flags[file_flags_index]
	return
}

// Return the name hash 2 at the the BET table index
func (archive *Archive) BetTableNameHash2Index(bet_table_index uint32) (name_hash_2 uint64, err error) {
	bet_table := &archive.bet_table

	if bet_table_index >= archive.BetTableCount() {
		err = fmt.Errorf("mpq: BET table name hash 2 index out of bounds: %d/%d", bet_table_index, archive.BetTableCount())
		return
	}

	name_hash_position := uint64(bet_table_index) * uint64(bet_table.header.BitTotal_NameHash2)
	name_hash_2, err = bet_table.name_hashes.Uint(name_hash_position, uint8(bet_table.header.BitCount_NameHash2))
	if err != nil {
		return
	}

	return
}
