package crypto

const (
	// mod is the largest prime that is less than 65536.
	adler_mod = 65521
	// nmax is the largest n such that
	// 255 * n * (n+1) / 2 + (n+1) * (mod-1) <= 2^32-1.
	// It is mentioned in RFC 1950 (search for "5552").
	adler_nmax = 5552
)

// Add p to the running checksum d.
func Adler32(d uint32, p []byte) uint32 {
	s1, s2 := uint32(d&0xffff), uint32(d>>16)
	for len(p) > 0 {
		var q []byte
		if len(p) > adler_nmax {
			p, q = p[:adler_nmax], p[adler_nmax:]
		}
		for len(p) >= 4 {
			s1 += uint32(p[0])
			s2 += s1
			s1 += uint32(p[1])
			s2 += s1
			s1 += uint32(p[2])
			s2 += s1
			s1 += uint32(p[3])
			s2 += s1
			p = p[4:]
		}
		for _, x := range p {
			s1 += uint32(x)
			s2 += s1
		}
		s1 %= adler_mod
		s2 %= adler_mod
		p = q
	}
	return s2<<16 | s1
}
