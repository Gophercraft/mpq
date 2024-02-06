package mpq

func (file *File) Close() (err error) {
	// close the underlying file descriptor to .MPQ
	err = file.file.Close()
	return
}
