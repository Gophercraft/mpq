package mpq

import "github.com/Gophercraft/mpq/info"

// Archive represents a *singular* .MPQ archive file
// To interact with multiple overlapping MPQs at once, refer to mpq.Set.
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
}
