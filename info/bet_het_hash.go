package info

// BetTableMergeHashValue reconstitutes the adjusted hash lookup value from name hashes 1 and 2
func BetTableMergeHashValue(bet_table_header *BetTableHeader, name_hash_1 uint8, name_hash_2 uint64) (name_hash_64 uint64) {
	return (uint64(name_hash_1) << uint64(bet_table_header.BitCount_NameHash2)) | name_hash_2
}

// Returns the Jenkins hash, but truncated to the required length, with the maximum allowed bit set to true.
// this allows the name_hash_1 value to always be non-zero
// in the HET table, a name_hash_1 zero would signify that there is no data
func HetTableLookupValue(het_table_header *HetTableHeader, in_hash_value uint64) (hash_value uint64) {
	var (
		and_mask_64 uint64
		or_mask_64  uint64
	)

	hash_bits_length := uint64(het_table_header.NameHashBitSize)

	if het_table_header.NameHashBitSize == 0x40 {
		// do nothing, hash is maximum quality
		and_mask_64 = 0xFFFFFFFFFFFFFFFF
	} else {
		// truncate the hash to the required length
		and_mask_64 = (1 << hash_bits_length) - 1
	}

	// Set the highest bit to 1, removing the possibility of the highest byte being zero
	or_mask_64 = 1 << (hash_bits_length - 1)

	hash_value = (in_hash_value & and_mask_64) | or_mask_64

	return
}
