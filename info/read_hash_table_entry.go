package info

import (
	"encoding/binary"
	"io"
)

// Read a HashTableEntry from a Reader in binary
func ReadHashTableEntry(reader io.Reader, hash_table_entry *HashTableEntry) error {
	var hash_entry_data [HashTableEntrySize]byte
	if _, err := io.ReadFull(reader, hash_entry_data[:]); err != nil {
		return err
	}
	// Read name hashes
	hash_entry_slice := hash_entry_data[:]
	hash_table_entry.Name1 = binary.LittleEndian.Uint32(hash_entry_slice[0:4])
	hash_entry_slice = hash_entry_slice[4:]
	hash_table_entry.Name2 = binary.LittleEndian.Uint32(hash_entry_slice[0:4])
	hash_entry_slice = hash_entry_slice[4:]
	// Read locale
	hash_table_entry.Locale = binary.LittleEndian.Uint16(hash_entry_slice[0:2])
	hash_entry_slice = hash_entry_slice[2:]
	// Read platform
	hash_table_entry.Platform = hash_entry_slice[0]
	hash_entry_slice = hash_entry_slice[1:]
	hash_table_entry.Reserved = hash_entry_slice[0]
	hash_entry_slice = hash_entry_slice[1:]
	// Read block index
	hash_table_entry.BlockIndex = binary.LittleEndian.Uint32(hash_entry_slice[0:4])
	// hash_entry_slice = hash_entry_slice[4:]

	return nil
}
