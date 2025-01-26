package info

import (
	"encoding/binary"
	"io"
)

const HetTableHeaderSize = 32

// The header occupies the very beginning of the encrypted data block immediately following the ExtTableHeader
type HetTableHeader struct {
	// Size of the entire HET table, including HetTableHeader (in bytes)
	TableSize uint32
	// Number of occupied entries in the HET table
	UsedEntryCount uint32
	// Total number of entries in the HET table
	TotalEntryCount uint32
	// Size of the name hash entry (in bits)
	NameHashBitSize uint32
	// Total size of file index (in bits)
	IndexSizeTotal uint32
	// Extra bits in the file index
	IndexSizeExtra uint32
	// Effective size of the file index (in bits)
	IndexSize uint32
	// Size of the block index subtable (in bytes)
	IndexTableSize uint32
}

func ReadHetTableHeader(reader io.Reader, het_table_header *HetTableHeader) (err error) {
	var data [HetTableHeaderSize]byte
	if _, err = io.ReadFull(reader, data[:]); err != nil {
		return
	}

	field := data[:]

	het_table_header.TableSize = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]

	het_table_header.UsedEntryCount = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]

	het_table_header.TotalEntryCount = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]

	het_table_header.NameHashBitSize = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]

	het_table_header.IndexSizeTotal = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]

	het_table_header.IndexSizeExtra = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]

	het_table_header.IndexSize = binary.LittleEndian.Uint32(field[:4])
	field = field[4:]

	het_table_header.IndexTableSize = binary.LittleEndian.Uint32(field[:4])

	return
}

// Returns the size of a single packed index, in bits
func HetTableIndexBitLength(het_table_header *HetTableHeader) (bit_length uint8) {
	max_value := het_table_header.UsedEntryCount

	for max_value > 0 {
		max_value >>= 1
		bit_length++
	}

	return
}

// Returns the number of bits needed to hold all packed indices in HET table
func HetTableIndicesBitLength(het_table_header *HetTableHeader) (bit_length uint64) {
	bit_length = uint64(HetTableIndexBitLength(het_table_header)) * uint64(het_table_header.TotalEntryCount)
	return
}

// Returns the number of bytes needed to hold all packed indices in HET table
func HetTableIndicesSize(het_table_header *HetTableHeader) (size uint64) {
	index_bit_length := HetTableIndicesBitLength(het_table_header)
	size = ((uint64(index_bit_length) + 7) / 8)
	return
}
