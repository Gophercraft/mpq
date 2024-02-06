package mpq

func (archive *Archive) Close() error {
	archive.block_table = nil
	archive.hash_table = nil
	archive.hi_block_table = nil
	return nil
}
