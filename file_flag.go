package mpq

import "github.com/Gophercraft/mpq/info"

// returns true if one or more bits in flag are also present in the file's flags
func (file *File) has_flag(flag info.FileFlag) bool {
	return file.flags&flag != 0
}
