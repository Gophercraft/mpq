package info

import (
	"crypto/md5"
)

const (
	MaxHeaderSize = 0x17FF
	Header0Size   = 24
	Header1Size   = 12
	Header2Size   = 24
	Header3Size   = 140
)

// MPQ file header
type Header struct {
	Header0
	Header1
	Header2
	Header3
}

type Header0 struct {
	// Size of MPQ archive
	// This field is deprecated in the Burning Crusade MoPaQ format, and the size of the archive
	// is calculated as the size from the beginning of the archive to the end of the hash table,
	// block table, or extended block table (whichever is largest).
	ArchiveSize uint32
	// 0 = Format 1 (up to The Burning Crusade)
	// 1 = Format 2 (The Burning Crusade and newer)
	// 2 = Format 3 (WoW - Cataclysm beta or newer)
	// 3 = Format 4 (WoW - Cataclysm beta or newer)
	Version uint16
	// Power of two exponent specifying the number of 512-byte disk sectors in each logical sector
	// in the archive. The size of each logical sector in the archive is 512 * 2^wBlockSize.
	SectorSize uint16
	// Offset to the beginning of the hash table, relative to the beginning of the archive.
	HashTablePos uint32
	// Offset to the beginning of the block table, relative to the beginning of the archive.
	BlockTablePos uint32
	// Number of entries in the hash table. Must be a power of two, and must be less than 2^16 for
	// the original MoPaQ format, or less than 2^20 for the Burning Crusade format.
	HashTableSize uint32
	// Number of entries in the block table
	BlockTableSize uint32
}

type Header1 struct {
	// Offset to the beginning of array of 16-bit high parts of file offsets.
	HiBlockTablePos64 uint64
	// High 16 bits of the hash table offset for large archives.
	HashTablePosHi uint16
	// High 16 bits of the block table offset for large archives.
	BlockTablePosHi uint16
}

type Header2 struct {
	// 64-bit version of the archive size
	ArchiveSize64 uint64
	// 64-bit position of the BET table
	BetTablePos64 uint64
	// 64-bit position of the HET table
	HetTablePos64 uint64
}

type Header3 struct {
	// Compressed size of the hash table
	HashTableSize64 uint64
	// Compressed size of the block table
	BlockTableSize64 uint64
	// Compressed size of the hi-block table
	HiBlockTableSize64 uint64
	// Compressed size of the HET block
	HetTableSize64 uint64
	// Compressed size of the BET block
	BetTableSize64 uint64
	// Size of raw data chunk to calculate MD5.
	// MD5 of each data chunk follows the raw file data.
	RawChunkSize uint32
	// Array of MD5's
	// MD5 of the block table before decryption
	MD5_BlockTable [md5.Size]byte
	// MD5 of the hash table before decryption
	MD5_HashTable [md5.Size]byte
	// MD5 of the hi-block table
	MD5_HiBlockTable [md5.Size]byte
	// MD5 of the BET table before decryption
	MD5_BetTable [md5.Size]byte
	// MD5 of the HET table before decryption
	MD5_HetTable [md5.Size]byte
	// MD5 of the MPQ header from signature to (including) MD5_HetTable
	MD5_MpqHeader [md5.Size]byte
}
