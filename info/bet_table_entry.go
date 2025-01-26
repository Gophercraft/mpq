package info

type BetTableEntry struct {
	// Offset of the beginning of the file, relative to the beginning of the archive.
	Position uint64
	// Compressed file size
	BlockSize uint64
	// Only valid if the block is a file; otherwise meaningless, and should be 0.
	// If the file is compressed, this is the size of the uncompressed file data.
	FileSize uint64
	// Flags for the file
	FlagsIndex uint32
}
