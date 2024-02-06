package compress

import (
	"bytes"
	"io"

	"github.com/JoshVarga/blast"
)

var (
	_ Decompressor = Decompress_PKWARE_DCL
)

func Decompress_PKWARE_DCL(in_bytes []byte) (out_bytes []byte, err error) {
	var read_closer io.ReadCloser
	read_closer, err = blast.NewReader(bytes.NewReader(in_bytes))
	if err != nil {
		return
	}
	out_bytes, err = io.ReadAll(read_closer)
	read_closer.Close()
	return
}
