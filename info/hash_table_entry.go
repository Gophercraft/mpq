package info

// The bytes needed to store a HashTableEntry in binary format
const HashTableEntrySize = 16

type HashTableEntry struct {
	// The hash of the file path, using method A.
	Name1 uint32
	// The hash of the file path, using method B.
	Name2 uint32
	// The language of the file. This is a Windows LANGID data type, and uses the same values.
	// 0 indicates the default language (American English), or that the file is language-neutral.
	Locale uint16
	// The platform the file is used for. 0 indicates the default platform.
	// No other values have been observed.
	Platform uint8
	Reserved uint8
	// If the hash table entry is valid, this is the index into the block table of the file.
	// Otherwise, one of the following two values:
	//  - FFFFFFFFh: Hash table entry is empty, and has always been empty.
	//               Terminates searches for a given file.
	//  - FFFFFFFEh: Hash table entry is empty, but was valid at some point (a deleted file).
	//               Does not terminate searches for a given file.
	BlockIndex uint32
}
