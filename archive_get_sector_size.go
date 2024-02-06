package mpq

func (archive *Archive) get_sector_size() (sector_size uint32, err error) {
	sector_size = 512 << uint32(archive.header.SectorSize)
	return
}
