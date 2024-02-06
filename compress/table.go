package compress

type Decompressor func(in_bytes []byte) (out_bytes []byte, err error)

type DecompressTableEntry struct {
	Type
	Decompressor
}

var DecompressTable = []DecompressTableEntry{
	{Bzip2, Decompress_bzip2},
	{PKWARE, Decompress_PKWARE_DCL},
	{Zlib, Decompress_zlib},
	{Huffman, Decompress_Huffman},
	{ADPCM_stereo, Decompress_ADPCM_stereo},
	{ADPCM_mono, Decompress_ADPCM_mono},
	{Sparse, Decompress_Sparse},
}
