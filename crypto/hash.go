package crypto

import "strings"

func HashString(input string, offset uint16) (hash uint32) {
	var seed1 uint32 = 0x7FED7FED
	var seed2 uint32 = 0xEEEEEEEE

	str := strings.ToUpper(input)

	for _, curChar := range str {
		value := block_encryption_table[offset+uint16(curChar)]
		seed1 = (value ^ (seed1 + seed2)) & 0xFFFFFFFF
		seed2 = (uint32(curChar) + seed1 + seed2 + (seed2 << 5) + 3) & 0xFFFFFFFF
	}

	return seed1
}
