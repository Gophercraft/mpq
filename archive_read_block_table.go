package mpq

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/Gophercraft/mpq/crypto"
	"github.com/Gophercraft/mpq/info"
)

const max_block_table_length = 50000000

// returns the absolute position of the Archive's block table
func (archive *Archive) get_block_table_position() (absolute_pos uint64, err error) {
	relative_pos := uint64(archive.header.BlockTablePos)
	// Add expanded bits if applicable
	if archive.header.Version >= 1 {
		relative_pos |= uint64(archive.header.BlockTablePosHi) << 32
	}

	absolute_pos = uint64(archive.archive_pos) + relative_pos
	return
}

// determines if a hi block table needs to be decoded
func (archive *Archive) contains_hi_block_table() bool {
	return archive.header.HiBlockTablePos64 != 0
}

// returns the absolute position of the Archive's hi-block table
func (archive *Archive) get_hi_block_table_position() (absolute_pos uint64, err error) {
	absolute_pos = uint64(archive.archive_pos) + archive.header.HiBlockTablePos64
	return
}

// returns the number of block table entries
func (archive *Archive) get_block_table_length() (length uint64, err error) {
	return uint64(archive.header.BlockTableSize), nil
}

// reads all block table entries into Archive
func (archive *Archive) read_block_table(file *os.File) (err error) {
	var block_table_position uint64
	var block_table_length uint64
	var hi_block_position uint64
	block_table_position, err = archive.get_block_table_position()
	if err != nil {
		return
	}
	if _, err = file.Seek(int64(block_table_position), io.SeekStart); err != nil {
		return
	}
	// Get length of block table
	block_table_length, err = archive.get_block_table_length()
	if err != nil {
		return
	}
	if block_table_length > max_block_table_length {
		err = fmt.Errorf("mpq: '%s' block table is too long to be loaded safely (%d entries)", archive.path, block_table_length)
		return
	}

	// Read encrypted block table from file
	encrypted_block_table := make([]byte, block_table_length*info.BlockTableEntrySize)
	if _, err = io.ReadFull(file, encrypted_block_table); err != nil {
		return
	}
	// Decrypt
	decrypt_seed := crypto.HashString("(block table)", info.HashFileKey)
	if err = crypto.Decrypt(decrypt_seed, encrypted_block_table); err != nil {
		return
	}
	decrypted_block_table_reader := bytes.NewReader(encrypted_block_table)

	// Allocate actual block table
	archive.block_table = make([]info.BlockTableEntry, block_table_length)

	// Load block table data
	for i := range archive.block_table {
		block_table_entry := &archive.block_table[i]
		if err = info.ReadBlockTableEntry(decrypted_block_table_reader, block_table_entry); err != nil {
			return
		}
	}

	// If applicable, also load the hi block offset table
	if archive.contains_hi_block_table() {
		hi_block_position, err = archive.get_hi_block_table_position()
		if err != nil {
			return
		}
		// seek to hi block offset data
		if _, err = file.Seek(int64(hi_block_position), io.SeekStart); err != nil {
			return
		}
		// read hi block offset data
		archive.hi_block_table = make([]uint16, block_table_length)
		if err = binary.Read(file, binary.LittleEndian, archive.hi_block_table); err != nil {
			return
		}
	}

	return
}
