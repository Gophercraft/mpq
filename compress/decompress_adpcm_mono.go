package compress

import "fmt"

var (
	_ Decompressor = Decompress_ADPCM_mono
)

func Decompress_ADPCM_mono(in_bytes []byte) (out_bytes []byte, err error) {
	// TODO: implement ADPCM
	err = fmt.Errorf("cannot decompress ADPCM mono")
	return
}
