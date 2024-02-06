package compress

import (
	"bytes"
	"compress/bzip2"
	"io"
)

var (
	_ Decompressor = Decompress_bzip2
)

func Decompress_bzip2(in_bytes []byte) (out_bytes []byte, err error) {
	return io.ReadAll(bzip2.NewReader(bytes.NewReader(in_bytes)))
}
