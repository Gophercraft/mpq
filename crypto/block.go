package crypto

import (
	"encoding/binary"
)

type block struct {
	key    uint32
	adjust uint32
}

func (d *block) BlockSize() int {
	return 4
}

func (b *block) Init(key uint32) {
	b.adjust = 0xEEEEEEEE
	b.key = key
}

func (d *block) Decrypt(dst, src []byte) {
	d.adjust += uint32(hash_source[(uint32(HashEncryptData)<<8)+(d.key&0xFF)])
	block := binary.LittleEndian.Uint32(src[:4]) ^ (d.key + d.adjust)
	d.key = ((^d.key << 0x15) + 0x11111111) | (d.key >> 0x0B)
	d.adjust = uint32(block) + d.adjust + (d.adjust << 5) + 3
	binary.LittleEndian.PutUint32(dst[:4], block)
}

func (d *block) Encrypt(dst, src []byte) {
	d.adjust += uint32(hash_source[(uint32(HashEncryptData)<<8)+(d.key&0xFF)])
	block := binary.LittleEndian.Uint32(src[:4])
	binary.LittleEndian.PutUint32(dst[:4], block^(d.key+d.adjust))
	d.key = ((^d.key << 0x15) + 0x11111111) | (d.key >> 0x0B)
	d.adjust = uint32(block) + d.adjust + (d.adjust << 5) + 3
}

// Creates a new encryption block
// block satisfies the cipher.Block
// interface
func NewBlock(key uint32) (b *block) {
	b = new(block)
	b.Init(key)
	return
}

// Decrypts MPQ block information in-situ using key
func Decrypt(key uint32, data []byte) {
	var b block
	b.Init(key)
	for len(data) >= 4 {
		b.Decrypt(data[:4], data[:4])
		data = data[4:]
	}
}
