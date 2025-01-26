package mpq

import "github.com/Gophercraft/mpq/info"

// The MPQ archive header
func (archive *Archive) Header() (header *info.Header) {
	return &archive.header
}
