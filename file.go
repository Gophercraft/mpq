package mpq

import (
	"io"
	"os"

	"github.com/Gophercraft/mpq/info"
)

// File represents a compressed file from within the MPQ
// For the MPQ Archive file itself, refer to Archive.
type File struct {
	// the Archive this was opened with
	archive *Archive
	// the fd to the underlying MPQ
	file *os.File
	// the flags found in the block table
	flags info.FileFlag
	// the locale (refer to info.HashTableEntry)
	locale uint16
	// the platform (refer to info.HashTableEntry)
	platform uint8
	// the size of the uncompressed file
	size uint64
	// the decryption key
	decryption_key uint32
	// the size of the compressed file
	compressed_size uint64
	// the location of the file's data in the MPQ archive
	archive_position uint64
	// the current sector
	sector_index int
	// sector offsets
	// potentially incorrect sidenote: it seems that the last sector offset is the "terminator" of the offset list
	sector_offsets []uint32
	// CRC32 checksums for each sector
	sector_checksums []uint32
	// the current sector reader
	sector_reader io.Reader
}
