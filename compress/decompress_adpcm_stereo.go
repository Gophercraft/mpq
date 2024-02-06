package compress

import "fmt"

var (
	_ Decompressor = Decompress_ADPCM_stereo
)

func Decompress_ADPCM_stereo(in_bytes []byte) (out_bytes []byte, err error) {
	// TODO: implement ADPCM
	err = fmt.Errorf("cannot decompress ADPCM stereo")
	return
}
