package mpq

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/Gophercraft/mpq/compress"
	"github.com/Gophercraft/mpq/crypto"
	"github.com/Gophercraft/mpq/info"
)

func (archive *Archive) contains_bet_table() (contains bool) {
	contains = archive.header.BetTablePos64 != 0
	return
}

func (archive *Archive) read_bet_table(file io.ReadSeeker) (err error) {
	// check if BET table exists
	if !archive.contains_bet_table() {
		return
	}
	// read entire BET table at once
	if _, err = file.Seek(archive.Position()+int64(archive.header.BetTablePos64), io.SeekStart); err != nil {
		return
	}
	bet_table_data := make([]byte, archive.header.BetTableSize64)
	if _, err = io.ReadFull(file, bet_table_data); err != nil {
		return
	}
	// get MD5 checksum of encrypted + compressed BET table
	bet_table_data_md5 := md5.Sum(bet_table_data[:])
	if bet_table_data_md5 != archive.header.MD5_BetTable {
		err = fmt.Errorf("mpq: invalid BET table MD5 checksum")
		return
	}
	// decrypt BET table
	decrypt_key := crypto.HashString("(block table)", crypto.HashEncryptKey)
	crypto.Decrypt(decrypt_key, bet_table_data[info.ExtTableHeaderSize:])
	var ext_table_header info.ExtTableHeader
	if err = info.ReadExtTableHeader(bytes.NewReader(bet_table_data[:info.ExtTableHeaderSize]), &ext_table_header); err != nil {
		return
	}
	if ext_table_header.Version != 1 {
		err = fmt.Errorf("mpq: BET table extended header: invalid version")
		return
	}
	if ext_table_header.Signature != info.BetTableSignature {
		err = fmt.Errorf("mpq: malformed BET table position")
		return
	}

	var bet_table []byte

	if ext_table_header.Size+info.ExtTableHeaderSize == uint32(len(bet_table_data)) {
		// no decompression required
		bet_table = bet_table_data[info.ExtTableHeaderSize:]
	} else {
		// decompression required
		bet_table, err = compress.Decompress2(bet_table_data[info.ExtTableHeaderSize:])
		if err != nil {
			return
		}
	}

	bet_table_reader := bytes.NewReader(bet_table)

	if err = info.ReadBetTableHeader(bet_table_reader, &archive.bet_table.header); err != nil {
		return
	}

	// read file flags
	archive.bet_table.file_flags = make([]info.FileFlag, archive.bet_table.header.FlagCount)
	file_flags := make([]byte, archive.bet_table.header.FlagCount*4)
	if _, err = io.ReadFull(bet_table_reader, file_flags); err != nil {
		return
	}
	for i := range archive.bet_table.file_flags {
		archive.bet_table.file_flags[i] = info.FileFlag(binary.LittleEndian.Uint32(file_flags[i*4 : (i+1)*4]))
	}

	// size of all table entries (in bits)
	table_entries_size := uint64(archive.bet_table.header.TableEntrySize) * uint64(archive.bet_table.header.EntryCount)

	table_entries := make([]byte, (table_entries_size+7)/8)
	if _, err = io.ReadFull(bet_table_reader, table_entries); err != nil {
		return
	}
	archive.bet_table.entries.Init(table_entries_size, table_entries)

	// read BET name hash2s
	name_hashes_size := uint64(archive.bet_table.header.BitTotal_NameHash2) * uint64(archive.bet_table.header.EntryCount)
	name_hashes := make([]byte, (name_hashes_size+7)/8)
	if _, err = io.ReadFull(bet_table_reader, name_hashes); err != nil {
		return
	}
	archive.bet_table.name_hashes.Init(name_hashes_size, name_hashes)

	if bet_table_reader.Len() != 0 {
		err = fmt.Errorf("mpq: error reading BET table: %d bytes unread", bet_table_reader.Len())
		return
	}

	// pBetTable->dwTableEntrySize * pBetHeader->dwEntryCount

	return
}
