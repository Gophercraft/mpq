package info

func LogicalSectorSize(header *Header) int {
	return int(uint32(512) << uint32(header.SectorSize))
}
