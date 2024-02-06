package info

const (
	HashTableIndex = 0x000
	HashNameA      = 0x100
	HashNameB      = 0x200
	HashFileKey    = 0x300
	HashKey2Mix    = 0x400
)

// Special block indices
const (
	HashTerminator uint32 = 0xFFFFFFFF
	HashRemoved    uint32 = 0xFFFFFFFE
)
