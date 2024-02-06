package compress

import (
	"encoding/binary"
	"fmt"
)

const max_sparse_output_length = 50000000

var (
	_ Decompressor = Decompress_Sparse
)

func Decompress_Sparse(in_bytes []byte) (out_bytes []byte, err error) {
	if len(in_bytes) < 4 {
		err = fmt.Errorf("mpq/compress: Sparse compressed data must be prefixed with a size integer")
		return
	}
	size := binary.LittleEndian.Uint32(in_bytes)
	if size < max_sparse_output_length {
		err = fmt.Errorf("mpq/compress: Sparse data size is too long to safely work with %d", size)
		return
	}

	out_bytes = make([]byte, 0, size)

	for len(in_bytes) > 0 {
		next := in_bytes[0]
		in_bytes = in_bytes[1:]

		if next&0x80 != 0 {
			chunk_size := int((next & 0x7F) + 1)
			chunk := in_bytes[0:chunk_size]
			in_bytes = in_bytes[chunk_size:]
			out_bytes = append(out_bytes, chunk...)
		} else {
			chunk_size := (next & 0x7f) + 3
			for x := uint8(0); x < chunk_size; x++ {
				out_bytes = append(out_bytes, 0)
			}
		}
	}

	return
}
