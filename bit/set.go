package bit

import (
	"fmt"
	"slices"
)

type Set struct {
	n    uint64
	data []byte
}

// Init initialize the set with n bits, and data storing your bits
func (set *Set) Init(n uint64, data []byte) {
	set.n = n
	set.data = data
}

// Len the length of the Set, in bits
func (set *Set) Len() (length uint64) {
	length = set.n
	return
}

func make_bit_flag(bit_offset byte) byte {
	return byte(1) << bit_offset
}

func (set *Set) Get(i uint64) (b bool, err error) {
	byte_offset := i / 8
	bit_offset := i % 8
	bit_flag := make_bit_flag(byte(bit_offset))

	if byte_offset > uint64(len(set.data)) {
		err = fmt.Errorf("bit: index %d out of range 0-%d", i, set.n)
		return
	}

	b = set.data[byte_offset]&bit_flag != 0
	return
}

func (set *Set) Set(i uint64, b bool) {
	byte_offset := i / 8
	bit_offset := i % 8
	bit_flag := make_bit_flag(byte(bit_offset))

	length := uint64(len(set.data))

	if length <= byte_offset {
		set.data = slices.Grow(set.data, int(byte_offset+1-length))
		set.n = i + 1
	}

	if b {
		set.data[byte_offset] |= bit_flag
	} else {
		set.data[byte_offset] &= ^bit_flag
	}
}

func (set *Set) Uint(first_bit uint64, bits uint8) (u uint64, err error) {
	if bits > 64 {
		err = fmt.Errorf("bit: too many bits to store in a 64-bit unsigned integer")
		return
	}

	var is_1 bool
	for i := uint64(0); i < uint64(bits); i++ {
		is_1, err = set.Get(first_bit + i)
		if err != nil {
			return
		}
		if is_1 {
			u |= (1 << i)
		}
	}

	return
}
