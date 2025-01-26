package mpq

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/Gophercraft/mpq/crypto"
	"github.com/Gophercraft/mpq/info"
)

func (file *File) init() (err error) {
	// The file can't be opened if it doesn't exist
	if !file.has_flag(info.FileExists) {
		err = fmt.Errorf("mpq: the file does not exist")
		return
	}

	// Open the MPQ Archive
	file.file, err = os.Open(file.archive.path)
	if err != nil {
		return
	}

	// Seek to the location of the encoded file
	if _, err = file.file.Seek(file.archive.archive_pos+int64(file.block_position), io.SeekStart); err != nil {
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
			crypto.Decrypt(file.decryption_key-1, sector_data)
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
			uint32(file.compressed_size),
		}
	}
	return
}
