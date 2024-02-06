package mpq

import (
	"bytes"
	"io"

	"github.com/Gophercraft/mpq/compress"
	"github.com/Gophercraft/mpq/crypto"
	"github.com/Gophercraft/mpq/info"
)

// create the next sector reader
func (file *File) read_next_sector() (sector io.Reader, err error) {
	// seek file to the location of the sector
	if _, err = file.file.Seek(int64(file.archive_position+uint64(file.sector_offsets[file.sector_index])), io.SeekStart); err != nil {
		return
	}
	// allocate space to hold raw sector data
	sector_data := make([]byte, file.sector_offsets[file.sector_index+1]-file.sector_offsets[file.sector_index])

	// read raw sector data
	if _, err = io.ReadFull(file.file, sector_data); err != nil {
		return
	}

	// decrypt sector if applicable
	if file.has_flag(info.FileEncrypted) {
		if err = crypto.Decrypt(file.decryption_key+uint32(file.sector_index), sector_data); err != nil {
			return
		}
	}

	// decompress data if applicable
	if file.has_flag(info.FileCompressMask) {
		// After MPQ v2, compression works a little differently
		if file.archive.header.Version >= 1 {
			sector_data, err = compress.Decompress2(sector_data)
		} else {
			sector_data, err = compress.Decompress1(sector_data)
		}
		if err != nil {
			return
		}
	}

	// return data as an io.Reader
	sector = bytes.NewReader(sector_data)
	return
}
