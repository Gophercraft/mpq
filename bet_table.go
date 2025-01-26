package mpq

import (
	"github.com/Gophercraft/mpq/bit"
	"github.com/Gophercraft/mpq/info"
)

type bet_table struct {
	header      info.BetTableHeader
	file_flags  []info.FileFlag
	entries     bit.Set
	name_hashes bit.Set
}
