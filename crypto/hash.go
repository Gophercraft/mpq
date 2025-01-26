package crypto

import (
	"encoding/binary"
	"hash"
	"strings"
)

// Hash types are permutations that are used to separate different types of hashes.
type HashType uint8

const (
	HashIndex HashType = iota
	// A and B hashes are used to reduce the rate of collisions
	HashNameA
	HashNameB
	// HashEncryptKey is used when generating decryption keys from strings
	HashEncryptKey
	// HashEncryptData is used internally when encrypting data
	HashEncryptData
)

type digest struct {
	result    uint32
	adjust    uint32
	hash_type HashType
}

func (d *digest) Reset() {
	d.result = 0x7FED7FED
	d.adjust = 0xEEEEEEEE
}

func (d *digest) Write(b []byte) (n int, err error) {
	for n = 0; n < len(b); n++ {
		d.result = (d.result + d.adjust) ^ uint32(hash_source[(uint32(d.hash_type)<<8)+uint32(b[n])])
		d.adjust += uint32(b[n]) + d.result + (d.adjust << 5) + 3
	}

	return
}

func (d *digest) BlockSize() int {
	return 1
}

func (d *digest) Size() int {
	return 4
}

func (d *digest) Sum32() (sum uint32) {
	sum = d.result
	return
}

func (d *digest) Sum(b []byte) []byte {
	var result [4]byte
	binary.LittleEndian.PutUint32(result[:], d.Sum32())
	return append(b, result[:]...)
}

// Creates a new Blizz hash digest. This digest implements the hash.Hash32 interface.
func NewHash(hash_type HashType) hash.Hash32 {
	d := new(digest)
	d.hash_type = hash_type
	d.Reset()
	return d
}

// Produces a hash for any given string. (case-insensitive)
func HashString(s string, hash_type HashType) uint32 {
	s = strings.ToUpper(s)
	h := NewHash(hash_type)
	h.Write([]byte(s))
	return h.Sum32()
}
