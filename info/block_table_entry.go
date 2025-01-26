package info

const BlockTableEntrySize = 16

type BlockTableEntry struct {
	// Offset of the beginning of the file, relative to the beginning of the archive.
	Position uint32
	// Compressed file size
	BlockSize uint32
	// Only valid if the block is a file; otherwise meaningless, and should be 0.
	// If the file is compressed, this is the size of the uncompressed file data.
	FileSize uint32
	// Flags for the file
	Flags FileFlag
}
