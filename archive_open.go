package mpq

import (
	"github.com/Gophercraft/mpq/crypto"
	"github.com/Gophercraft/mpq/info"
)

// Open opens a [File] contained within the [Archive]
func (archive *Archive) Open(name string) (file *File, err error) {
	// calculate decryption key

	if archive.contains_het_table() {
		// we can take advantage of the version 4 lookup method
		file, err = archive.open_ext(name)
		return
	}

	// use the classic method
	// Lookup name in hash table
	var hash_table_entry *info.HashTableEntry
	_, hash_table_entry, err = archive.hash_lookup_name(name)
	if err != nil {
		return
	}

	// Find block table entry
	var block_table_entry *info.BlockTableEntry
	block_table_entry, err = archive.BlockTableIndex(hash_table_entry.BlockIndex)
	if err != nil {
		return
	}

	// Get the position of the encoded file
	block_position := uint64(block_table_entry.Position)
	// The file size might be higher if it's encoded with more than 32 bits
	// Apply an additional 16 bits from the hi block table
	if archive.ContainsHiBlockTable() {
		var hi_block_position uint16
		hi_block_position, err = archive.HiBlockTableIndex(hash_table_entry.BlockIndex)
		if err != nil {
			return
		}
		block_position |= uint64(hi_block_position) << 32
	}
	// Begin constructing file object
	file = new(File)
	file.name = name
	file.archive = archive
	file.block_position = block_position
	file.flags = block_table_entry.Flags
	file.compressed_size = uint64(block_table_entry.BlockSize)
	file.size = uint64(block_table_entry.FileSize)
	file.decryption_key = crypto.HashString(name, crypto.HashEncryptKey)

	// Fix decryption key
	if file.has_flag(info.FileEncrypted) && file.has_flag(info.FileFixKey) {
		file.decryption_key = (file.decryption_key + block_table_entry.Position) ^ block_table_entry.BlockSize
	}

	err = file.init()
	if err != nil {
		return
	}

	// Populate these attributes from the hash table entry
	file.locale = hash_table_entry.Locale
	file.platform = hash_table_entry.Platform

	return
}
