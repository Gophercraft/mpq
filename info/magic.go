package info

var (
	UserDataSignature   = [4]byte{'M', 'P', 'Q', 0x1B}
	HeaderDataSignature = [4]byte{'M', 'P', 'Q', 0x1A}

	HetTableSignature = [4]byte{'H', 'E', 'T', 0x1A}
	BetTableSignature = [4]byte{'B', 'E', 'T', 0x1A}
)
