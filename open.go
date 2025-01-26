package mpq

import (
	"os"
	"path/filepath"
)

// Open will open the MoPaQ [Archive] specified by name
func Open(name string) (archive *Archive, err error) {
	archive = new(Archive)
	// Get fully qualified path to file
	archive.path, err = filepath.Abs(name)
	if err != nil {
		// Just use regular path if this errors
		archive.path = name
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

	// read HET table
	if err = archive.read_het_table(file); err != nil {
		return
	}

	// read BET table
	if err = archive.read_bet_table(file); err != nil {
		return
	}

	// close fd
	err = file.Close()
	if err != nil {
		return
	}

	return
}
