package mpq

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Gophercraft/mpq/crypto"
	"github.com/Gophercraft/mpq/info"
)

const max_hash_table_length = 50000000

// returns the absolute position of the Archive's hash table
func (archive *Archive) get_hash_table_position() (absolute_pos uint64, err error) {
	absolute_pos = uint64(archive.archive_pos) + info.HashTablePos(&archive.header)
	return
}

func (archive *Archive) get_hash_table_length() (length uint64, err error) {
	return uint64(archive.header.HashTableSize), nil
}

// read all hash table entries into Archive
func (archive *Archive) read_hash_table(file io.ReadSeeker) (err error) {
	// seek to the start of the hash table
	var hash_table_position uint64
	var hash_table_length uint64
	hash_table_position, err = archive.get_hash_table_position()
	if err != nil {
		return
	}
	if _, err = file.Seek(int64(hash_table_position), io.SeekStart); err != nil {
		return
	}

	// Allocate hash table
	hash_table_length, err = archive.get_hash_table_length()
	if err != nil {
		return
	}
	if hash_table_length > max_hash_table_length {
		err = fmt.Errorf("mpq: '%s' hash table is too long (%d entries)", archive.path, hash_table_length)
		return
	}
	archive.hash_table = make([]info.HashTableEntry, hash_table_length)

	// allocate temporary space to hold the encrypted data
	hash_table_data := make([]byte, info.HashTableEntrySize*hash_table_length)

	// read encrypted hash table
	if _, err = io.ReadFull(file, hash_table_data); err != nil {
		return
	}

	// decrypt hash table
	decrypt_seed := crypto.HashString("(hash table)", crypto.HashEncryptKey)
	crypto.Decrypt(decrypt_seed, hash_table_data)

	// read hash table
	hash_table_reader := bytes.NewReader(hash_table_data)
	for i := range archive.hash_table {
		hash_table_entry := &archive.hash_table[i]
		if err = info.ReadHashTableEntry(hash_table_reader, hash_table_entry); err != nil {
			return
		}
	}

	return
}
