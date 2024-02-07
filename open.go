package mpq

import (
	"os"
	"path/filepath"
)

func Open(path string) (archive *Archive, err error) {
	archive = new(Archive)
	// Get fully qualified path to file
	archive.path, err = filepath.Abs(path)
	if err != nil {
		// Just use regular path if this errors
		archive.path = path
		err = nil
	}

	// open fd to archive
	file, err := os.Open(archive.path)
	if err != nil {
		return
	}

	// read user data and header
	err = archive.read_header(file)
	if err != nil {
		return
	}

	// read hash table
	err = archive.read_hash_table(file)
	if err != nil {
		return
	}

	// read block table and hi block table
	err = archive.read_block_table(file)
	if err != nil {
		return
	}

	// close fd
	err = file.Close()
	if err != nil {
		return
	}

	// pre-calculate information about reading files
	archive.sector_size = int(uint32(512) << uint32(archive.header.SectorSize))

	return
}
