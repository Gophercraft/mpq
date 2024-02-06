package info

const BlockTableEntrySize = 16

type BlockTableEntry struct {
	// Offset of the beginning of the file, relative to the beginning of the archive.
	FilePos uint32
	// Compressed file size
	CompressedSize uint32
	// Only valid if the block is a file; otherwise meaningless, and should be 0.
	// If the file is compressed, this is the size of the uncompressed file data.
	DecompressedSize uint32
	// Flags for the file
	Flag FileFlag
}
