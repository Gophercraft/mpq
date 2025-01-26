package mpq

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"

	"github.com/Gophercraft/mpq/compress"
	"github.com/Gophercraft/mpq/crypto"
	"github.com/Gophercraft/mpq/info"
)

func (archive *Archive) read_het_table(file io.ReadSeeker) (err error) {
	// check if HET table exists
	if archive.header.HetTablePos64 == 0 {
		return
	}
	// read entire HET table at once
	if _, err = file.Seek(archive.Position()+int64(archive.header.HetTablePos64), io.SeekStart); err != nil {
		return
	}
	het_table_data := make([]byte, archive.header.HetTableSize64)
	if _, err = io.ReadFull(file, het_table_data); err != nil {
		return
	}
	// get MD5 checksum of encrypted + compressed HET table
	het_table_data_md5 := md5.Sum(het_table_data[:])
	if het_table_data_md5 != archive.header.MD5_HetTable {
		err = fmt.Errorf("mpq: invalid HET table MD5 checksum")
		return
	}
	// decrypt HET table
	decrypt_key := crypto.HashString("(hash table)", crypto.HashEncryptKey)
	crypto.Decrypt(decrypt_key, het_table_data[info.ExtTableHeaderSize:])
	var ext_table_header info.ExtTableHeader
	if err = info.ReadExtTableHeader(bytes.NewReader(het_table_data[:info.ExtTableHeaderSize]), &ext_table_header); err != nil {
		return
	}
	if ext_table_header.Signature != info.HetTableSignature {
		err = fmt.Errorf("mpq: malformed HET table position")
		return
	}

	var het_table []byte

	if ext_table_header.Size+info.ExtTableHeaderSize == uint32(len(het_table_data)) {
		// no decompression required
		het_table = het_table_data[info.ExtTableHeaderSize:]
	} else {
		// decompression required
		het_table, err = compress.Decompress2(het_table_data[info.ExtTableHeaderSize:])
		if err != nil {
			return
		}
	}

	het_table_reader := bytes.NewReader(het_table)

	// Read table header
	if err = info.ReadHetTableHeader(het_table_reader, &archive.het_table.header); err != nil {
		return
	}

	// Read hashes
	name_hashes := make([]byte, archive.het_table.header.TotalEntryCount)
	if _, err = io.ReadFull(het_table_reader, name_hashes[:]); err != nil {
		return
	}
	archive.het_table.name_hashes = name_hashes

	// Read indices
	bet_indices := make([]byte, info.HetTableIndicesSize(&archive.het_table.header))
	if _, err = io.ReadFull(het_table_reader, bet_indices); err != nil {
		return
	}

	archive.het_table.bet_indices.Init(info.HetTableIndicesBitLength(&archive.het_table.header), bet_indices)

	if het_table_reader.Len() != 0 {
		err = fmt.Errorf("mpq: HET table read failure")
		return
	}

	return
}
