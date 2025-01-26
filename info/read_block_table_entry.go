package info

import (
	"encoding/binary"
	"io"
)

// Read a BlockTableEntry from a Reader in binary
func ReadBlockTableEntry(reader io.Reader, block_table_entry *BlockTableEntry) error {
	var block_table_entry_data [BlockTableEntrySize]byte
	if _, err := io.ReadFull(reader, block_table_entry_data[:]); err != nil {
		return err
	}

	entry_slice := block_table_entry_data[:]

	// Read file pos
	block_table_entry.Position = binary.LittleEndian.Uint32(entry_slice[:4])
	entry_slice = entry_slice[4:]

	// Read compressed size
	block_table_entry.BlockSize = binary.LittleEndian.Uint32(entry_slice[:4])
	entry_slice = entry_slice[4:]

	// Read uncompressed size
	block_table_entry.FileSize = binary.LittleEndian.Uint32(entry_slice[:4])
	entry_slice = entry_slice[4:]

	// Read file flags
	block_table_entry.Flags = FileFlag(binary.LittleEndian.Uint32(entry_slice[:4]))

	return nil
}
