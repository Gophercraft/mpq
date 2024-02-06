package compress

import (
	"bytes"
	"compress/zlib"
	"io"
)

var (
	_ Decompressor = Decompress_zlib
)

func Decompress_zlib(in_bytes []byte) (out_bytes []byte, err error) {
	var read_closer io.ReadCloser
	read_closer, err = zlib.NewReader(bytes.NewReader(in_bytes))
	if err != nil {
		return
	}
	out_bytes, err = io.ReadAll(read_closer)
	if err != nil {
		return
	}
	read_closer.Close()
	return
}
