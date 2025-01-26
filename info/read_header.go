package info

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
)

func ReadHeader(reader io.Reader, header *Header) (err error) {
	// read header size
	var header_size_data [4]byte
	var header_size uint32
	if _, err = io.ReadFull(reader, header_size_data[:]); err != nil {
		return err
	}
	header_size = binary.LittleEndian.Uint32(header_size_data[:])
	if header_size > MaxHeaderSize {
		return fmt.Errorf("info: header size is too large (%d)", header_size)
	}

	// read header data
	header_data := make([]byte, header_size-8)
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

	if header.Version >= 3 {
		// check for header integrity
		integrity_check_data := append(append(HeaderDataSignature[:], header_size_data[:]...), header_data[:len(header_data)-md5.Size]...)
		mpq_header_md5 := md5.Sum(integrity_check_data)

		if mpq_header_md5 != header.MD5_MpqHeader {
			err = fmt.Errorf("info: MPQ header MD5 integrity check failed (hash in header is %s, hash of header data is %s)", hex.EncodeToString(header.MD5_MpqHeader[:]), hex.EncodeToString(mpq_header_md5[:]))
			return
		}
	}

	header_reader = nil
	return
}
