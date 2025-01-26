package mpq

import (
	"fmt"

	"github.com/Gophercraft/mpq/crypto"
	"github.com/Gophercraft/mpq/info"
)

// Open a file using the newer (extended) lookup system
func (archive *Archive) open_ext(name string) (file *File, err error) {
	var (
		bet_table_index uint32
		bet_table_entry info.BetTableEntry
		file_flags      info.FileFlag
	)
	bet_table_index, err = archive.ext_table_lookup(name)
	if err != nil {
		return
	}
	if err = archive.BetTableEntryIndex(bet_table_index, &bet_table_entry); err != nil {
		return
	}
	file_flags, err = archive.BetTableFileFlagsIndex(bet_table_entry.FlagsIndex)
	if err != nil {
		return
	}

	// Begin constructing file object
	file = new(File)
	file.name = name
	file.archive = archive
	file.block_position = bet_table_entry.Position
	file.flags = file_flags
	file.compressed_size = bet_table_entry.BlockSize
	file.size = bet_table_entry.FileSize
	file.decryption_key = crypto.HashString(name, crypto.HashEncryptKey)

	if file.has_flag(info.FileEncrypted) && file.has_flag(info.FileFixKey) {
		err = fmt.Errorf("mpq: unknown if extended tables support fix key encryption")
		return
	}

	err = file.init()
	return
}
