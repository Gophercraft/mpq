package mpq

import "github.com/Gophercraft/mpq/info"

// An Archive provides access to a MoPaQ archive file
type Archive struct {
	// The fully qualified path pointing to the .MPQ file
	path string
	// The full size of the file
	file_size int64
	// The position relative to start of file where the MPQ archive begins
	archive_pos int64
	// The MPQ header
	header info.Header
	// The MPQ user data
	user_data info.UserData
	// The hash table
	hash_table []info.HashTableEntry
	// The block table
	block_table []info.BlockTableEntry
	// The table holding extended file offset bits
	hi_block_table []uint16
	// The HET table
	het_table het_table
	// The BET table
	bet_table bet_table
}
