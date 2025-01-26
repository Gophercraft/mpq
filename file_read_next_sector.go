package mpq

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Gophercraft/mpq/compress"
	"github.com/Gophercraft/mpq/crypto"
	"github.com/Gophercraft/mpq/info"
)

// create the next sector reader
func (file *File) read_next_sector() (sector io.Reader, err error) {
	// if the file is empty, succeed immediately with an EOF
	if file.size == 0 {
		err = io.EOF
		return
	}

	// seek file to the location of the sector
	if _, err = file.file.Seek(file.Position()+int64(file.sector_offsets[file.sector_index]), io.SeekStart); err != nil {
		err = fmt.Errorf("mpq: could not seek to next sector: %w", err)
		return
	}

	// calculate the length to hold a sector
	sector_length := file.sector_offsets[file.sector_index+1] - file.sector_offsets[file.sector_index]
	if uint64(sector_length) > file.size {
		err = fmt.Errorf("mpq: calculated sector offset is far too large to be safely read: %d", sector_length)
		return
	}

	// allocate space to hold raw sector data
	sector_data := make([]byte, sector_length)

	// read raw sector data
	if _, err = io.ReadFull(file.file, sector_data); err != nil {
		err = fmt.Errorf("mpq: could read raw sector data (%d): %w", len(sector_data), err)
		return
	}

	// decrypt sector if applicable
	if file.has_flag(info.FileEncrypted) {
		crypto.Decrypt(file.decryption_key+uint32(file.sector_index), sector_data)
	}

	// decompress data if applicable
	if file.has_flag(info.FileCompressMask) {
		// 12340 enUS\Data\common-2.MPQ : Character\BROKEN\Female\BrokenFemale0062-00.anim
		// this file has compress mask, but is not actually compressed at all

		// Sometimes the last sector can be inferred to be uncompressed
		// because it completes the file by addition
		// or its length is equal to the MPQ header's sector size
		needs_decompression_applied := !(uint64(len(sector_data))+file.bytes_read == file.size || len(sector_data) == int(info.LogicalSectorSize(&file.archive.header)))

		if needs_decompression_applied {
			var decompressed_data []byte
			// After MPQ v2, compression works a little differently
			if file.archive.header.Version >= 1 {
				decompressed_data, err = compress.Decompress2(sector_data)
			} else {
				decompressed_data, err = compress.Decompress1(sector_data)
			}

			if err != nil {
				return
			}
			sector_data = decompressed_data
		}
	}

	// return sector as an io.Reader
	sector = bytes.NewReader(sector_data)
	return
}
