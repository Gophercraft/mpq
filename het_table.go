package mpq

import (
	"github.com/Gophercraft/mpq/bit"
	"github.com/Gophercraft/mpq/info"
)

type het_table struct {
	header      info.HetTableHeader
	name_hashes []byte
	bet_indices bit.Set
}
