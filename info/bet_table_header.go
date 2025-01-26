package info

import (
	"encoding/binary"
	"io"
)

const BetTableHeaderSize = 76

type BetTableHeader struct {
	// Size of the entire BET table, including the header (in bytes)
	TableSize uint32
	// Number of entries in the BET table. Must match HET_TABLE_HEADER::dwEntryCount
	EntryCount uint32
	Unknown08  uint32
	// Size of one table entry (in bits)
	TableEntrySize uint32
	// Bit index of the file position (within the entry record)
	BitIndex_FilePos uint32
	// Bit index of the file size (within the entry record)
	BitIndex_FileSize uint32
	// Bit index of the compressed size (within the entry record)
	BitIndex_CompressedSize uint32
	// Bit index of the flag index (within the entry record)
	BitIndex_FlagIndex uint32
	// Bit index of the ??? (within the entry record)
	BitIndex_Unknown uint32
	// Bit size of file position (in the entry record)
	BitCount_FilePos uint32
	// Bit size of file size (in the entry record)
	BitCount_FileSize uint32
	// Bit size of compressed file size (in the entry record)
	BitCount_CompressedSize uint32
	// Bit size of flags index (in the entry record)
	BitCount_FlagIndex uint32
	// Bit size of ??? (in the entry record)
	BitCount_Unknown uint32
	// Total bit size of the NameHash2
	BitTotal_NameHash2 uint32
	// Extra bits in the NameHash2
	BitExtra_NameHash2 uint32
	// Effective size of NameHash2 (in bits)
	BitCount_NameHash2 uint32
	// Size of NameHash2 table, in bytes
	NameHashArraySize uint32
	// Number of flags in the following array
	FlagCount uint32
}

func ReadBetTableHeader(reader io.Reader, bet_table_header *BetTableHeader) (err error) {
	var data [BetTableHeaderSize]byte
	_, err = io.ReadFull(reader, data[:])
	if err != nil {
		return
	}

	field := data[:]

	bet_table_header.TableSize = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.EntryCount = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.Unknown08 = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.TableEntrySize = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.BitIndex_FilePos = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.BitIndex_FileSize = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.BitIndex_CompressedSize = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.BitIndex_FlagIndex = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.BitIndex_Unknown = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.BitCount_FilePos = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.BitCount_FileSize = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.BitCount_CompressedSize = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.BitCount_FlagIndex = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.BitCount_Unknown = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.BitTotal_NameHash2 = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.BitExtra_NameHash2 = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.BitCount_NameHash2 = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.NameHashArraySize = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]
	bet_table_header.FlagCount = binary.LittleEndian.Uint32(field[:4])

	return
}
