package mpq

import (
	"github.com/Gophercraft/mpq/info"
)

// return the number of sectors used to store a file
func (file *File) get_sector_count() (sector_count uint32, err error) {
	sector_size := uint32(file.archive.sector_size)

	if file.has_flag(info.FileSingleUnit) {
		// a single unit file only ever contains one sector.
		sector_count = 1
	} else {
		// divide the file size (rounded up a sector) by the size of a sector
		// to get the logical amount of sectors
		sector_count = (uint32(file.size) + sector_size - 1) / sector_size
	}

	return
}
