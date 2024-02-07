package mpq

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/Gophercraft/mpq/crypto"
	"github.com/Gophercraft/mpq/info"
)

// Opens a File contained within the archive
func (archive *Archive) Open(path string) (file *File, err error) {
	// Lookup file in hash table
	var hash_table_entry *info.HashTableEntry
	_, hash_table_entry, err = archive.hash_lookup_name(path)
	if err != nil {
		return
	}
	// calculate decryption key
	file_key := crypto.HashString(path, info.HashFileKey)

	// Find block table entry
	block_index := int(hash_table_entry.BlockIndex)
	// Ensure that block is within bounds
	if block_index >= len(archive.block_table) {
		err = fmt.Errorf("mpq: the hash table entry of '%s' points to a block table entry that is out of bounds (%d)", path, block_index)
		return
	}
	block_table_entry := &archive.block_table[block_index]

	// Get the position of the encoded file
	file_position := uint64(block_table_entry.FilePos)
	// The file size might be higher if it's encoded with more than 32 bits
	// Apply an additional 16 bits from the hi block table
	if archive.contains_hi_block_table() {
		file_position |= uint64(archive.hi_block_table[block_index]) << 32
	}
	// Make relative position into absolute
	file_position = uint64(archive.archive_pos) + file_position
	// Begin constructing file object
	file = new(File)
	file.path = path
	file.archive = archive
	file.archive_position = file_position
	file.flags = block_table_entry.Flag
	file.compressed_size = uint64(block_table_entry.CompressedSize)
	file.size = uint64(block_table_entry.DecompressedSize)
	file.decryption_key = file_key

	// Fix decryption key
	if file.has_flag(info.FileEncrypted) && file.has_flag(info.FileFixKey) {
		file.decryption_key = (file.decryption_key + block_table_entry.FilePos) ^ block_table_entry.DecompressedSize
	}

	// Populate these attributes from the hash table entry
	file.locale = hash_table_entry.Locale
	file.platform = hash_table_entry.Platform

	// Open the MPQ Archive
	file.file, err = os.Open(archive.path)
	if err != nil {
		return
	}

	// Seek to the location of the encoded file
	if _, err = file.file.Seek(int64(file_position), io.SeekStart); err != nil {
		return
	}

	// some files are implicitly single-unit (speech-enUS.MPQ : BloodElfAmbience.mp3)
	implicit_single_unit := !file.has_flag(info.FileCompress) && file.compressed_size == file.size

	if !file.has_flag(info.FileSingleUnit) && !implicit_single_unit {
		// calculate number of sectors
		var sector_count uint32
		sector_count, err = file.get_sector_count()
		if err != nil {
			return
		}
		file.sector_count = int(sector_count)

		// we add 1 to store the upper limit of the final sector
		num_sector_offsets := file.sector_count + 1
		// there is a "sector" to store the checksum table
		if file.has_flag(info.FileSectorCRC) {
			num_sector_offsets++
		}
		// Read offsets for multi-sector files
		sector_data := make([]byte, num_sector_offsets*4)
		if _, err = io.ReadFull(file.file, sector_data); err != nil {
			return
		}
		// Decrypt sector data
		if file.has_flag(info.FileEncrypted) {
			if err = crypto.Decrypt(file.decryption_key-1, sector_data); err != nil {
				return
			}
		}

		// Decode sector data
		file.sector_offsets = make([]uint32, num_sector_offsets)
		for i := range file.sector_offsets {
			file.sector_offsets[i] = binary.LittleEndian.Uint32(sector_data[i*4 : (i+1)*4])
		}
	} else {
		file.sector_count = 1
		// Since this is flagged as a single unit file, the offsets can be inferred as a simple slice
		file.sector_offsets = []uint32{
			0,
			block_table_entry.CompressedSize,
		}
	}

	return
}
