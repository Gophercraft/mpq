package compress

import (
	"bytes"
	"io"

	"github.com/ulikunitz/xz/lzma"
)

var (
	_ Decompressor = Decompress_LZMA
)

func Decompress_LZMA(in_bytes []byte) (out_bytes []byte, err error) {
	var read_closer *lzma.Reader
	read_closer, err = lzma.NewReader(bytes.NewReader(in_bytes))
	if err != nil {
		return
	}
	out_bytes, err = io.ReadAll(read_closer)
	return
}
