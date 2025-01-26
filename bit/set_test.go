package bit

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	var set Set

	check_bit := func(i uint64, value bool) {
		v, err := set.Get(i)
		if err != nil {
			panic(err)
		}
		if v != value {
			t.Fatal("bit", v, "does not match expected", value)
		}
	}

	check := func(i uint64, bits uint8, value uint64) {
		u, err := set.Uint(i, bits)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("bits", i, bits, value)
		if u != value {
			t.Fatal("value does not match expected", i, bits, value, "was", u)
		}
	}

	set.Init(24, []byte{0xFE, 0xFE, 0xFE})

	check_bit(0, false)
	check_bit(1, true)
	check_bit(2, true)
	check_bit(3, true)

	check(0, 1, 0)
	check(1, 1, 1)
	check(2, 1, 1)
	check(3, 1, 1)
	check(4, 1, 1)
	check(5, 1, 1)
	check(6, 1, 1)
	check(7, 1, 1)

	check(8, 1, 0)
	check(9, 1, 1)
	check(10, 1, 1)
	check(11, 1, 1)
	check(12, 1, 1)
	check(13, 1, 1)
	check(14, 1, 1)
	check(15, 1, 1)

	check(16, 1, 0)
	check(17, 1, 1)
	check(18, 1, 1)
	check(19, 1, 1)
	check(20, 1, 1)
	check(21, 1, 1)
	check(22, 1, 1)
	check(23, 1, 1)

	// check 2-bit numbers (endianness)
	check(0, 2, 2)
	check(2, 2, 3)
	check(4, 2, 3)
	check(6, 2, 3)

	check(8, 2, 2)
	check(10, 2, 3)
	check(12, 2, 3)
	check(14, 2, 3)

	check(16, 2, 2)
	check(18, 2, 3)
	check(20, 2, 3)
	check(22, 2, 3)
}
