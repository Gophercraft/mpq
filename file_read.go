package mpq

import (
	"fmt"
	"io"
)

func (file *File) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return
	}
	for {
		// If there is no reader, get the next sector
		if file.sector_reader == nil {
			file.sector_reader, err = file.read_next_sector()
			if err != nil {
				return
			}
		}

		// read from the current sector
		n, err = file.sector_reader.Read(b)

		file.bytes_read += uint64(n)

		if err == io.EOF {
			// if the current sector has been read completely
			// move to next sector
			file.sector_reader = nil
			file.sector_index++

			// If this is the last sector offset, then we're truly out of data to read
			if file.sector_index == file.sector_count {
				// if we've reached the end of the file,
				// the number of bytes read must be equal to the uncompressed file size
				// if not, a serious error has occurred
				if file.bytes_read != file.size {
					err = fmt.Errorf("mpq: reached end of file '%s', but # of bytes read (%d) mismatches the size of the file (%d)", file.path, file.bytes_read, file.size)
				}
				return
			}

			// otherwise, we have more data
			// so no error
			err = nil

			// Nop if no bytes were read, act as if nothing happened!
			if n == 0 {
				continue
			}
			// return n, nil
			return
		}

		return
	}
}
