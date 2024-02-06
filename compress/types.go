package compress

type Type uint32

const (
	Huffman      Type = 0x01       // Huffmann compression (used on WAVE files only)
	Zlib         Type = 0x02       // ZLIB compression
	PKWARE       Type = 0x08       // PKWARE DCL compression
	Bzip2        Type = 0x10       // BZIP2 compression (added in Warcraft III)
	Sparse       Type = 0x20       // Sparse compression (added in Starcraft 2)
	ADPCM_mono   Type = 0x40       // IMA ADPCM compression (mono)
	ADPCM_stereo Type = 0x80       // IMA ADPCM compression (stereo)
	LZMA         Type = 0x12       // LZMA compression. Added in Starcraft 2. This value is NOT a combination of flags.
	NextSame     Type = 0xFFFFFFFF // Same compression
)
