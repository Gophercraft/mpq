package compress

import "fmt"

func Decompress1(in_bytes []byte) (out_bytes []byte, err error) {
	const valid_mask uint8 = 0xFF

	// Get applied compression types and decrement data length
	compression_mask_1 := Type(in_bytes[0] & valid_mask)
	compression_mask_2 := compression_mask_1
	in_bytes = in_bytes[1:]

	// This compression function doesn't support LZMA
	if compression_mask_1 == LZMA {
		err = fmt.Errorf("mpq/compress: Decompress1 does not support LZMA (0x%02x)", compression_mask_1)
		return
	}

	table := DecompressTable

	var compress_count int

	// Parse the compression mask
	for i := 0; i < len(DecompressTable); i++ {
		// If the mask agrees, insert the compression function to the array
		if compression_mask_1&table[i].Type != 0 {
			compression_mask_2 &= ^(table[i].Type)
			compress_count++
		}
	}

	// If at least one of the compressions remaing unknown, return an error
	if compress_count == 0 || compression_mask_2 != 0 {
		err = fmt.Errorf("mpq/compress: unsupported compression bit in mask")
		return
	}

	// Apply all decompressions
	for i := 0; i < len(table); i++ {
		// Perform the (next) decompression
		if compression_mask_1&table[i].Type != 0 {
			in_bytes, err = table[i].Decompressor(in_bytes)
			if err != nil {
				return
			}
		}
	}

	out_bytes = in_bytes

	return
}

func Decompress2(in_bytes []byte) (out_bytes []byte, err error) {
	// Get the compression methods
	compression_method := Type(in_bytes[0])
	in_bytes = in_bytes[1:]
	var decompressors [2]Decompressor

	// We only recognize a fixed set of compression methods
	switch compression_method {
	case Zlib:
		decompressors[0] = Decompress_zlib
	case PKWARE:
		decompressors[0] = Decompress_PKWARE_DCL
	case Bzip2:
		decompressors[0] = Decompress_bzip2
	case LZMA:
		decompressors[0] = Decompress_LZMA
	case Sparse:
		decompressors[0] = Decompress_Sparse
	case Sparse | Zlib:
		decompressors[0] = Decompress_zlib
		decompressors[1] = Decompress_Sparse
	case Sparse | Bzip2:
		decompressors[0] = Decompress_bzip2
		decompressors[1] = Decompress_Sparse
	case ADPCM_mono | Huffman:
		decompressors[0] = Decompress_Huffman
		decompressors[1] = Decompress_ADPCM_mono
	case ADPCM_stereo | Huffman:
		decompressors[0] = Decompress_Huffman
		decompressors[1] = Decompress_ADPCM_mono
	default:
		return nil, fmt.Errorf("mpq/compress: file corrupted (invalid compression method 0x%02x)", uint8(compression_method))
	}

	// Apply the first decompression method
	out_bytes, err = decompressors[0](in_bytes)
	if err != nil {
		err = fmt.Errorf("mpq: first decompressor, %s", err)
		return
	}

	// Apply the second decompression method, if any
	if decompressors[1] != nil {
		out_bytes, err = decompressors[1](out_bytes)
		if err != nil {
			err = fmt.Errorf("mpq: second decompressor, %s", err)
			return
		}
	}

	return
}
