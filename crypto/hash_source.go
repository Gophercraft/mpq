package crypto

var hash_source [256 * 5]uint32

// setup the lookup table before the program runs
func init() {
	seed := uint32(0x100001)
	for i := uint32(0); i < 256; i++ {
		for j := uint32(0); j < 5; j++ {
			seed = (seed*0x7D + 3) % 0x2AAAAB
			rand1 := seed & 0xFFFF
			seed = (seed*0x7D + 3) % 0x2AAAAB
			rand2 := seed & 0xFFFF
			hash_source[(j<<8)+i] = (rand1 << 16) | rand2
		}
	}
}
