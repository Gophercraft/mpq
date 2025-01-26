package mpq

func (file *File) Size() uint64 {
	return file.size
}

// Position returns the absolute position of the file inside the [Archive].
func (file *File) Position() (pos int64) {
	return file.archive.archive_pos + int64(file.block_position)
}

// BlockSize returns the block size, or compressed size of the file
func (file *File) BlockSize() (bs uint64) {
	bs = file.compressed_size
	return
}
