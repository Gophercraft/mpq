package compress

import "fmt"

var (
	_ Decompressor = Decompress_Huffman
)

func Decompress_Huffman(in_bytes []byte) (out_bytes []byte, err error) {
	// TODO: implement Huffman decoding in separate file
	// https://github.com/ladislav-zezula/StormLib/blob/8978bdc8214f4ec3543ac11962fc094c3b0803b2/src/huffman/huff.cpp
	err = fmt.Errorf("cannot decompress huffman coding")
	return
}
