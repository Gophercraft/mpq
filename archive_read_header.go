package mpq

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/Gophercraft/mpq/info"
)

const header_alignment = 512

func (archive *Archive) read_header(file *os.File) (err error) {
	// seek to end of file (gets file size)
	archive.file_size, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return
	}

	// seek to start of file
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	var current_offset int64

	// attempt to find the header
	var magic_bytes [4]byte
	for {
		// get current offset
		current_offset, err = file.Seek(0, io.SeekCurrent)
		if err != nil {
			return
		}

		if _, err = io.ReadFull(file, magic_bytes[:]); err != nil {
			err = fmt.Errorf("failed to read magic bytes (current offset %d): %s", current_offset, err)
			return
		}

		// check for signatures of MPQ info
		switch {
		case bytes.Equal(magic_bytes[:], info.UserDataSignature):
			// read user data
			err = info.ReadUserData(file, &archive.user_data)
			if err != nil {
				return
			}
			new_offset := current_offset + int64(archive.user_data.HeaderOffset)
			// seek to header info
			_, err = file.Seek(new_offset, io.SeekStart)
			if err != nil {
				err = fmt.Errorf("failed to seek to position of MPQ header pointed to by user data: %s", err)
				return
			}
		case bytes.Equal(magic_bytes[:], info.HeaderDataSignature):
			// read header info
			err = info.ReadHeader(file, &archive.header)
			if err != nil {
				err = fmt.Errorf("failed to read MPQ header info: %s", err)
				return
			}
			// mark current position as the absolute beginning of the MPQ archive
			archive.archive_pos = current_offset
			return
		default:
			// nothing found (seek current + header_alignment)
			if _, err = file.Seek(header_alignment, io.SeekCurrent); err != nil {
				err = fmt.Errorf("failed to find MPQ header: %s", err)
				return
			}
		}
	}

}
