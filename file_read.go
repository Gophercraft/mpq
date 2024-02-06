package mpq

import (
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
		if err == io.EOF {
			// if the current sector has been read completely
			// move to next sector
			file.sector_reader = nil
			file.sector_index++

			// If this is the last sector offset, then we're truly out of data to read
			if file.sector_index == len(file.sector_offsets)-1 {
				// return n, io.EOF
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
