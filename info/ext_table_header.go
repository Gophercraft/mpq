package info

import (
	"encoding/binary"
	"io"
)

const (
	ExtTableHeaderSize = 12
)

type ExtTableHeader struct {
	// 'HET\x1A' or 'BET\x1A'
	Signature [4]byte
	// Version. Seems to be always 1
	Version uint32
	// Size of the contained table
	Size uint32
}

func ReadExtTableHeader(reader io.Reader, ext_table_header *ExtTableHeader) (err error) {
	var data [ExtTableHeaderSize]byte
	if _, err = io.ReadFull(reader, data[:]); err != nil {
		return
	}
	field := data[:]
	copy(ext_table_header.Signature[:], field[:4])
	field = field[4:]

	ext_table_header.Version = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]

	ext_table_header.Size = binary.LittleEndian.Uint32(field[:4])

	return
}
