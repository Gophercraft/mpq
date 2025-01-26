package cmd

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/Gophercraft/mpq"
	"github.com/Gophercraft/mpq/crypto"
	"github.com/Gophercraft/mpq/info"
	"github.com/Gophercraft/mpq/jenkins"
	"github.com/spf13/cobra"
)

type registry_blizz_hashes_key struct {
	NameA uint32
	NameB uint32
}

type hash_registry struct {
	archive        *mpq.Archive
	blizz_hashes   map[registry_blizz_hashes_key]string
	jenkins_hashes map[uint64]string
}

func new_hash_registry(archive *mpq.Archive) (registry *hash_registry) {
	registry = new(hash_registry)
	registry.archive = archive
	registry.blizz_hashes = make(map[registry_blizz_hashes_key]string)
	registry.jenkins_hashes = make(map[uint64]string)
	return
}

func (registry *hash_registry) LookupBlizz(name_a, name_b uint32) (name string, err error) {
	var ok bool
	name, ok = registry.blizz_hashes[registry_blizz_hashes_key{NameA: name_a, NameB: name_b}]
	if !ok {
		err = fmt.Errorf("could not find blizz hash in registry")
	}
	return
}

func (registry *hash_registry) LookupJenkins(name64 uint64) (name string, err error) {
	var ok bool
	name, ok = registry.jenkins_hashes[name64]
	if !ok {
		err = fmt.Errorf("could not find jenkins hash in registry")
	}
	return
}

func (registry *hash_registry) Add(name string) {
	// register blizz hashes
	var blizz_hash_key registry_blizz_hashes_key
	blizz_hash_key.NameA = crypto.HashString(name, crypto.HashNameA)
	blizz_hash_key.NameB = crypto.HashString(name, crypto.HashNameB)
	registry.blizz_hashes[blizz_hash_key] = name

	if registry.archive.Header().HetTablePos64 != 0 {
		// register jenkins hash
		jenkins_hash_key := info.HetTableLookupValue(registry.archive.HetTableHeader(), jenkins.Hash64([]byte(strings.ToLower(name))))
		registry.jenkins_hashes[jenkins_hash_key] = name
	}
}

// maps block table index to hash table index
type inverse_block_map map[uint32]uint32

var info_cmd = &cobra.Command{
	Use:   "info [mpq file]",
	Short: "read info about an MPQ archive",
	Run:   run_info_cmd,
}

func init() {
	f := info_cmd.Flags()
	f.Bool("hash-table", false, "display all the information contained within the hash table")
	f.Bool("block-table", false, "display all the information contained within the block table")
	f.Bool("het-table", false, "display all the information contained within the HET table")
	f.Bool("bet-table", false, "display all the information contained within the BET table")
	root_cmd.AddCommand(info_cmd)
}

var (
	g_tabwriter *tabwriter.Writer
	titles      int
)

func title(s string) {
	titles++
	g_tabwriter.Flush()
	if titles != 1 {
		fmt.Println()
	}

	c := make([]byte, len(s))
	for i := range c {
		c[i] = '-'
	}
	fmt.Println(s)
	fmt.Println(string(c))
}

func attrib(name string, args ...any) {
	g_tabwriter.Write([]byte(name))
	g_tabwriter.Write([]byte("\t"))
	fmt.Fprintln(g_tabwriter, args...)
}

func format_locale(locale uint16) string {
	var locales = map[uint16]string{
		0: "Universal",
	}

	name, ok := locales[locale]
	if ok {
		name = " (" + name + ")"
	}

	return fmt.Sprintf("0x%04X%s", locale, name)
}

func format_platform(platform uint8) string {
	return fmt.Sprintf("%d", platform)
}

func run_info_cmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		os.Exit(1)
	}

	flags := cmd.Flags()
	show_hash_table, err := flags.GetBool("hash-table")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	show_block_table, err := flags.GetBool("block-table")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	show_het_table, err := flags.GetBool("het-table")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	show_bet_table, err := flags.GetBool("bet-table")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	g_tabwriter = tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	archive_path := args[0]
	archive_path_abs, _ := filepath.Abs(archive_path)
	archive, err := mpq.Open(archive_path)
	if err != nil {
		fmt.Println("cannot open archive:", err)
		return
	}

	header := archive.Header()

	het_table_available := header.HetTablePos64 != 0
	bet_table_available := header.BetTablePos64 != 0

	title(archive_path_abs)

	attrib("Archive begins at", archive.Position())

	attrib("MPQ format version:", header.Version)

	if header.Version < 1 {
		attrib("Archive size:", header.ArchiveSize)
	} else {
		attrib("Archive size (deprecated):", header.ArchiveSize)
	}
	if header.Version >= 2 {
		attrib("Archive size (64-bit):", header.ArchiveSize64)
	}

	attrib("Sector size exponent:", header.SectorSize, fmt.Sprintf("(logical sector size computes as: %d)", info.LogicalSectorSize(header)))

	attrib("Hash table position:", info.HashTablePos(header), "(", header.HashTablePos, "(low)", header.HashTablePosHi, "(hi)", ")")

	attrib("Block table position:", info.BlockTablePos(header), "(", header.BlockTablePos, "(low)", header.BlockTablePosHi, "(hi)", ")")

	attrib("Hash table entries:", header.HashTableSize)
	if header.Version >= 3 {
		attrib("Hash table size (64-bit):", header.HashTableSize64)
	}

	attrib("Block table entries:", header.BlockTableSize)
	if header.Version >= 3 {
		attrib("Block table size (64-bit):", header.BlockTableSize64)
	}

	if header.Version >= 1 {
		attrib("Hi-block table position (64-bit):", header.HiBlockTablePos64)
	}

	if header.Version >= 2 {
		attrib("BET table position (64-bit):", header.BetTablePos64)
		attrib("HET table position (64-bit):", header.HetTablePos64)
	}

	if header.Version >= 3 {
		attrib("Compressed hash table size (64-bit):", header.HashTableSize64)
		attrib("Compressed block table table size (64-bit):", header.BlockTableSize64)
		attrib("Compressed hi-block table size (64-bit):", header.HiBlockTableSize64)
		attrib("Compressed HET table size (64-bit):", header.HetTableSize64)
		attrib("Compressed BET table size (64-bit):", header.BetTableSize64)
		attrib("Raw chunk size:", header.RawChunkSize)
		attrib("MD5 of block table before decryption:", hex.EncodeToString(header.MD5_BlockTable[:]))
		attrib("MD5 of hash table before decryption:", hex.EncodeToString(header.MD5_HashTable[:]))
		attrib("MD5 of hi-block table before decryption:", hex.EncodeToString(header.MD5_HiBlockTable[:]))
		attrib("MD5 of BET table before decryption:", hex.EncodeToString(header.MD5_BetTable[:]))
		attrib("MD5 of HET table before decryption:", hex.EncodeToString(header.MD5_HetTable[:]))
		attrib("MD5 of MPQ header from to end of MD5_HetTable:", hex.EncodeToString(header.MD5_MpqHeader[:]))
	}

	registry := new_hash_registry(archive)

	registry.Add("(listfile)")
	registry.Add("(attributes)")
	registry.Add("(signature)")
	registry.Add("(user data)")

	list, err := archive.List()
	if err != nil {
		attrib("List file does not exist, error: ", err)
	} else {
		i := 0
		for list.Next() {
			name := list.Path()
			registry.Add(name)
			i++
		}
		list.Close()
		attrib("List file exists", i, "entries")
	}

	if het_table_available {
		title("HET table header")

		het_table_header := archive.HetTableHeader()
		attrib("Size of entire HET table:", het_table_header.TableSize)
		attrib("Number of use entries:", het_table_header.UsedEntryCount)
		attrib("Total number of entries:", het_table_header.TotalEntryCount)
		attrib("Size of a name hash (in bits):", het_table_header.NameHashBitSize)
		attrib("Effective size of the file index:", het_table_header.IndexSize)
		attrib("Extra bits file index:", het_table_header.IndexSizeExtra)
		attrib("Total size of file index:", het_table_header.IndexSizeTotal)
		attrib("Block index subtable size:", het_table_header.IndexTableSize)
	}

	if bet_table_available {
		bet_table_header := archive.BetTableHeader()

		title("BET table header")

		attrib("Size of the entire BET table:", bet_table_header.TableSize)
		attrib("Number of entries in the BET table:", bet_table_header.EntryCount)
		attrib("Unknown:", bet_table_header.Unknown08)
		attrib("Table entry size (in bits):", bet_table_header.TableEntrySize)
		attrib("Bit index of the file position:", bet_table_header.BitIndex_FilePos)
		attrib("Bit index of the file size:", bet_table_header.BitIndex_FileSize)
		attrib("Bit index of the file compressed size:", bet_table_header.BitIndex_CompressedSize)
		attrib("Bit index of the flag index:", bet_table_header.BitIndex_FlagIndex)
		attrib("Bit index of unknown field:", bet_table_header.BitIndex_Unknown)

		attrib("Bit count of the file position:", bet_table_header.BitCount_FilePos)
		attrib("Bit count of the file size:", bet_table_header.BitCount_FileSize)
		attrib("Bit count of the file compressed size:", bet_table_header.BitCount_CompressedSize)
		attrib("Bit count of the flag index:", bet_table_header.BitCount_FlagIndex)
		attrib("Bit count of unknown field:", bet_table_header.BitCount_Unknown)

		attrib("Total bit size of name hash 2:", bet_table_header.BitTotal_NameHash2)
		attrib("Extra bits in name hash 2:", bet_table_header.BitExtra_NameHash2)
		attrib("Effective size of name hash 2:", bet_table_header.BitCount_NameHash2)

		attrib("Size of name hash 2 table in bytes:", bet_table_header.NameHashArraySize)
		attrib("Number of file flags in the lookup table:", bet_table_header.FlagCount)

		// 		// Size of the entire BET table, including the header (in bytes)
		// TableSize
		// // Number of entries in the BET table. Must match HET_TABLE_HEADER::dwEntryCount
		// EntryCount uint32
		// Unknown08  uint32
		// // Size of one table entry (in bits)
		// TableEntrySize uint32
		// // Bit index of the file position (within the entry record)
		// BitIndex_FilePos uint32
		// // Bit index of the file size (within the entry record)
		// BitIndex_FileSize uint32
		// // Bit index of the compressed size (within the entry record)
		// BitIndex_CompressedSize uint32
		// // Bit index of the flag index (within the entry record)
		// BitIndex_FlagIndex uint32
		// // Bit index of the ??? (within the entry record)
		// BitIndex_Unknown uint32
		// // Bit size of file position (in the entry record)
		// BitCount_FilePos uint32
		// // Bit size of file size (in the entry record)
		// BitCount_FileSize uint32
		// // Bit size of compressed file size (in the entry record)
		// BitCount_CompressedSize uint32
		// // Bit size of flags index (in the entry record)
		// BitCount_FlagIndex uint32
		// // Bit size of ??? (in the entry record)
		// BitCount_Unknown uint32
		// // Total bit size of the NameHash2
		// BitTotal_NameHash2 uint32
		// // Extra bits in the NameHash2
		// BitExtra_NameHash2 uint32
		// // Effective size of NameHash2 (in bits)
		// BitCount_NameHash2 uint32
		// // Size of NameHash2 table, in bytes
		// NameHashArraySize uint32
		// // Number of flags in the following array
		// FlagCount uint32
	}

	// show hash table if requested
	if show_hash_table {
		title("Hash table")
		attrib("Number of entries:", archive.HashTableCount())
		attrib("Position (absolute):", archive.Position()+int64(info.HashTablePos(header)))
		for i := range archive.HashTableCount() {

			hash_table_entry, err := archive.HashTableIndex(i)
			if err != nil {
				attrib(" Error indexing hash table", err)
			} else {
				// if hash_table_entry.BlockIndex == 0xFFFFFFFF {
				// attrib(" Deleted entry (no block index)")
				// } else {
				if hash_table_entry.BlockIndex != 0xFFFFFFFF {
					if i > 0 {
						attrib(" --")
					}
					attrib(" Index:", i)
					attrib(" Name hash A:", fmt.Sprintf("%08x", hash_table_entry.Name1))
					attrib(" Name hash B:", fmt.Sprintf("%08x", hash_table_entry.Name2))
					corresponding_filename, lookup_err := registry.LookupBlizz(hash_table_entry.Name1, hash_table_entry.Name2)
					if lookup_err == nil {
						attrib(" Hashes correspond to a filename:", corresponding_filename)
					} else {
						attrib(" Entry hashes are not known!")
					}

					attrib(" Locale", format_locale(hash_table_entry.Locale))
					attrib(" Platform", hash_table_entry.Platform)

					attrib(" Block table index", hash_table_entry.BlockIndex)
					// }
				}
			}
		}
	}

	inverse_block_table := make(inverse_block_map)
	for i := range archive.HashTableCount() {
		hash_table_entry, err := archive.HashTableIndex(i)
		if err == nil {
			inverse_block_table[hash_table_entry.BlockIndex] = i
		}
	}

	if show_block_table {
		flag_usage := make(map[info.FileFlag]uint32)

		title("Block table")
		attrib("Number of entries:", archive.BlockTableCount())
		attrib("Position (absolute):", archive.Position()+int64(info.BlockTablePos(header)))
		for i := range archive.BlockTableCount() {
			if i > 0 {
				attrib(" --")
			}

			block_table_entry, err := archive.BlockTableIndex(i)
			if err != nil {
				attrib(" Error indexing block table", err)
			} else {
				attrib(" Index:", i)
				hash_table_index, found := inverse_block_table[i]
				if found {
					hash_table_entry, err := archive.HashTableIndex(hash_table_index)
					if err != nil {
						panic(err)
					}
					attrib(" In hash table: ")
					filename, lookup_error := registry.LookupBlizz(hash_table_entry.Name1, hash_table_entry.Name2)
					if lookup_error == nil {
						attrib("  Filename:", filename)
					}
					attrib("  Locale:", format_locale(hash_table_entry.Locale))
					attrib("  Platform:", format_platform(hash_table_entry.Platform))
				}
				block_position := uint64(block_table_entry.Position)
				var hi_block_position uint16
				if archive.ContainsHiBlockTable() {
					hi_block_position, err = archive.HiBlockTableIndex(i)
					if err != nil {
						attrib("Error getting hi-block position", err)
					} else {
						block_position |= uint64(hi_block_position) << 32
					}
				}

				attrib(" File position:", block_position)
				attrib("  Absolute file position:", archive.Position()+int64(block_position))
				if archive.ContainsHiBlockTable() {
					attrib("  Uses hi-block table?", "yes")
				}

				attrib(" Compressed size:", block_table_entry.BlockSize)
				attrib(" Decompressed size:", block_table_entry.FileSize)
				attrib(" Flags:", block_table_entry.Flags)
				if block_table_entry.Flags&info.FileEncrypted != 0 {
					flag_usage[info.FileEncrypted]++
				}
				if block_table_entry.Flags&info.FileFixKey != 0 {
					flag_usage[info.FileFixKey]++
				}
			}
		}

		attrib("Stats:")
		attrib("% of fix-key obfuscated:", float64(flag_usage[info.FileEncrypted])/float64(archive.BlockTableCount())*100.0)
		attrib("% of encrypted:", float64(flag_usage[info.FileFixKey])/float64(archive.BlockTableCount())*100.0)
	}

	if het_table_available && show_het_table {
		title("HET table")

		attrib("Number of entries", archive.HetTableCount())

		for i := range archive.HetTableCount() {
			name_hash_1, err := archive.HetTableNameHash1Index(i)
			if err != nil {
				attrib(" Error indexing name hash", err)
			}
			if name_hash_1 != 0x00 {
				if i > 0 {
					attrib(" --")
				}
				attrib(" Index:", i)

				bet_table_index, err := archive.HetTableIndexBetTableIndex(i)
				if err != nil {
					attrib(" Error indexing BET table indices", err)
				}
				attrib(" Name hash 1", fmt.Sprintf("%02x", name_hash_1))
				attrib(" BET table index", bet_table_index)

				name_hash_2, err := archive.BetTableNameHash2Index(bet_table_index)
				if err == nil {
					name_hash_64 := info.BetTableMergeHashValue(archive.BetTableHeader(), name_hash_1, name_hash_2)
					attrib(" Complete name hash:", fmt.Sprintf("%016x", name_hash_64))
					corresponding_filename, lookup_err := registry.LookupJenkins(name_hash_64)
					if lookup_err == nil {
						attrib(" Corresponds to filename:", corresponding_filename)
					}
				} else {
					panic(err)
				}
			}
		}
	}

	name_hash1s_bet_indices := make(map[uint32]uint8)
	if het_table_available {
		// associate block indexes with name hash1s
		for i := range archive.HetTableCount() {
			name_hash_1_value, err := archive.HetTableNameHash1Index(i)
			if err != nil {
				panic(err)
			}
			bet_table_index, err := archive.HetTableIndexBetTableIndex(i)
			if err != nil {
				panic(err)
			}
			name_hash1s_bet_indices[bet_table_index] = name_hash_1_value
		}
	}

	if bet_table_available && show_bet_table {
		flag_usage := make(map[info.FileFlag]uint32)

		title("BET table")

		attrib("Number of entries", archive.BetTableCount())
		for i := range archive.BetTableCount() {
			if i > 0 {
				attrib(" --")
			}
			attrib("Index:", i)
			name_hash_2, err := archive.BetTableNameHash2Index(i)
			if err == nil {
				attrib("Name hash 2:", fmt.Sprintf("%014x", name_hash_2))
				name_hash_1, ok := name_hash1s_bet_indices[i]
				if ok {
					name_hash := info.BetTableMergeHashValue(archive.BetTableHeader(), name_hash_1, name_hash_2)
					attrib("Complete name hash:", fmt.Sprintf("%016x", name_hash))
					corresponding_name, err := registry.LookupJenkins(name_hash)
					if err == nil {
						attrib("Corresponding name:", corresponding_name)
					}
				}
			}
			var bet_table_entry info.BetTableEntry
			err = archive.BetTableEntryIndex(i, &bet_table_entry)
			if err == nil {
				block_position := uint64(bet_table_entry.Position)
				attrib(" File position:", block_position)
				attrib("  Absolute file position:", archive.Position()+int64(block_position))
				if archive.ContainsHiBlockTable() {
					attrib("  Uses hi-block table?", "yes")
				}

				attrib(" Compressed size:", bet_table_entry.BlockSize)
				attrib(" Decompressed size:", bet_table_entry.FileSize)
				attrib(" Flags index:", bet_table_entry.FlagsIndex)

				bet_table_flags, err := archive.BetTableFileFlagsIndex(bet_table_entry.FlagsIndex)
				if err == nil {
					attrib(" Flags:", bet_table_flags)
					if bet_table_flags&info.FileEncrypted != 0 {
						flag_usage[info.FileEncrypted]++
					}
					if bet_table_flags&info.FileFixKey != 0 {
						flag_usage[info.FileFixKey]++
					}
				}

				attrib("Stats:")
				attrib("% of fix-key obfuscated:", float64(flag_usage[info.FileEncrypted])/float64(archive.BetTableCount())*100.0)
				attrib("% of encrypted:", float64(flag_usage[info.FileFixKey])/float64(archive.BetTableCount())*100.0)
			}
		}
	}

	archive.Close()

	g_tabwriter.Flush()

	if !show_hash_table {
		fmt.Println()
		fmt.Fprintln(os.Stderr, "hint: you can use --hash-table to inspect the contents of the archive's hash table.")
	}

	if !show_block_table {
		fmt.Println()
		fmt.Fprintln(os.Stderr, "hint: you can use --block-table to inspect the contents of the archive's block (file) table.")
	}

	if !show_het_table && het_table_available {
		fmt.Println()
		fmt.Fprintln(os.Stderr, "hint: you can use --het-table to inspect the contents of the archive's HET table (extended hash table)")
	}

	if !show_bet_table && bet_table_available {
		fmt.Println()
		fmt.Fprintln(os.Stderr, "hint: you can use --bet-table to inspect the contents of the archive's BET table (extended block table)")
	}
}
