package jenkins

func rot(x, k uint32) uint32 {
	return (((x) << (k)) | ((x) >> (32 - (k))))
}

func mix(ia, ib, ic uint32) (a, b, c uint32) {
	a, b, c = ia, ib, ic

	a -= c
	a ^= rot(c, 4)
	c += b
	b -= a
	b ^= rot(a, 6)
	a += c
	c -= b
	c ^= rot(b, 8)
	b += a
	a -= c
	a ^= rot(c, 16)
	c += b
	b -= a
	b ^= rot(a, 19)
	a += c
	c -= b
	c ^= rot(b, 4)
	b += a

	return
}

func final(ia, ib, ic uint32) (a, b, c uint32) {
	a, b, c = ia, ib, ic

	c ^= b
	c -= rot(b, 14)
	a ^= c
	a -= rot(c, 11)
	b ^= a
	b -= rot(a, 25)
	c ^= b
	c -= rot(b, 16)
	a ^= c
	a -= rot(c, 4)
	b ^= a
	b -= rot(a, 14)
	c ^= b
	c -= rot(b, 24)
	return
}

// returns a 64-bit hash. Equivalent to SBaseCommon.cpp:ULONGLONG HashStringJenkins(const char * szFileName)
// only without string conversion
// https://github.com/ladislav-zezula/StormLib/blob/d1bf5e1c71af432ccc24371a3611dc33edb7361f/src/SBaseCommon.cpp#L370
func Hash64(key []byte) (result uint64) {
	var (
		a, b, c uint32
		length  = uint32(len(key))
	)

	// 2 = MPQ initval
	c = 0xdeadbeef + length + 2
	b = c
	a = b
	// 1 = MPQ initval
	c += 1

	k := key

	// all but the last block: affect some 32 bits of (a,b,c)
	for length > 12 {
		a += uint32(k[0])
		a += (uint32(k[1])) << 8
		a += (uint32(k[2])) << 16
		a += (uint32(k[3])) << 24
		b += uint32(k[4])
		b += (uint32(k[5])) << 8
		b += (uint32(k[6])) << 16
		b += (uint32(k[7])) << 24
		c += uint32(k[8])
		c += (uint32(k[9])) << 8
		c += (uint32(k[10])) << 16
		c += (uint32(k[11])) << 24
		a, b, c = mix(a, b, c)
		length -= 12
		k = k[12:]
	}

	// last block: affect all 32 bits of (c)
	// all the case statements fall through
	switch length {
	case 12:
		c += (uint32(k[11])) << 24
		fallthrough
	case 11:
		c += (uint32(k[10])) << 16
		fallthrough
	case 10:
		c += (uint32(k[9])) << 8
		fallthrough
	case 9:
		c += uint32(k[8])
		fallthrough
	case 8:
		b += (uint32(k[7])) << 24
		fallthrough
	case 7:
		b += (uint32(k[6])) << 16
		fallthrough
	case 6:
		b += (uint32(k[5])) << 8
		fallthrough
	case 5:
		b += uint32(k[4])
		fallthrough
	case 4:
		a += (uint32(k[3])) << 24
		fallthrough
	case 3:
		a += (uint32(k[2])) << 16
		fallthrough
	case 2:
		a += (uint32(k[1])) << 8
		fallthrough
	case 1:
		a += uint32(k[0])
	case 0:
		result = uint64(b)<<32 | uint64(c)
		return // zero length strings require no mixing
	}

	_, b, c = final(a, b, c)
	result = uint64(b)<<32 | uint64(c)
	return
}
