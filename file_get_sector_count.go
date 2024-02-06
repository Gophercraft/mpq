package mpq

import "github.com/Gophercraft/mpq/info"

// return the number of sectors used to store a file
func (file *File) get_sector_count() (sector_count uint32, err error) {
	var sector_size uint32
	sector_size, err = file.archive.get_sector_size()
	if err != nil {
		return 0, err
	}

	if file.has_flag(info.FileSingleUnit) {
		sector_count = 1
	} else {
		sector_count = (uint32(file.size)+sector_size-1)/sector_size + 1
	}
	return
}
