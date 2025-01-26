package info

func HashTablePos(header *Header) (pos uint64) {
	pos = uint64(header.HashTablePos) | uint64(header.HashTablePosHi)<<32
	return
}
