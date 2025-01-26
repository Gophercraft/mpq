package mpq

import (
	"io"
	"os"

	"github.com/Gophercraft/mpq/info"
)

// File represents a handle to a file contained within an [Archive]
type File struct {
	// the path that was used to open the File
	name string
	// the Archive this was opened with
	archive *Archive
	// the fd to the underlying MPQ
	file *os.File
	// the locale (refer to info.HashTableEntry)
	locale uint16
	// the platform (refer to info.HashTableEntry)
	platform uint8
	// the decryption key
	decryption_key uint32
	// the location of the file's block relative to the start of the archive
	block_position uint64
	// the size of the uncompressed file
	size uint64
	// the size of the compressed file block
	compressed_size uint64
	// the flags found in the block table
	flags info.FileFlag
	// the number of sectors
	sector_count int
	// the current sector
	sector_index int
	// sector offsets
	// potentially incorrect sidenote: it seems that the last sector offset [sector_count-1] is the "terminator" of the offset list
	// not a marker for the beginning of a sector like the previous elements
	sector_offsets []uint32
	// number of bytes read
	bytes_read uint64
	// the current sector reader
	sector_reader io.Reader
}
