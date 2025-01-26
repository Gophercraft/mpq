package info

func BlockTablePos(header *Header) (pos uint64) {
	pos = uint64(header.BlockTablePos) | uint64(header.BlockTablePosHi)<<32
	return
}
