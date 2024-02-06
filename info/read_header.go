package info

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func ReadHeader(reader io.Reader, header *Header) (err error) {
	// read header size
	var header_size uint32
	if err = binary.Read(reader, binary.LittleEndian, &header_size); err != nil {
		return err
	}
	if header_size > MaxHeaderSize {
		return fmt.Errorf("info: header size is too large (%d)", header_size)
	}

	// read header data
	header_data := make([]byte, header_size)
	if _, err = io.ReadFull(reader, header_data); err != nil {
		return
	}
	header_reader := bytes.NewReader(header_data)
	// read first segment of header (version 1.0)
	if err = binary.Read(header_reader, binary.LittleEndian, &header.Header0); err != nil {
		return err
	}

	switch header.Version {
	case 0:
		// nothing further to read
	case 1:
		// read MPQ v2 header
		if err = binary.Read(header_reader, binary.LittleEndian, &header.Header1); err != nil {
			return err
		}
	case 2:
		// read MPQ v3 header
		if err = binary.Read(header_reader, binary.LittleEndian, &header.Header1); err != nil {
			return err
		}
		if err = binary.Read(header_reader, binary.LittleEndian, &header.Header2); err != nil {
			return err
		}
	case 3:
		// read MPQ v4 header
		if err = binary.Read(header_reader, binary.LittleEndian, &header.Header1); err != nil {
			return err
		}
		if err = binary.Read(header_reader, binary.LittleEndian, &header.Header2); err != nil {
			return err
		}
		if err = binary.Read(header_reader, binary.LittleEndian, &header.Header3); err != nil {
			return err
		}
	default:
		err = fmt.Errorf("info: error while reading Header: unknown format version (%d)", header.Version)
	}

	header_reader = nil
	return
}
