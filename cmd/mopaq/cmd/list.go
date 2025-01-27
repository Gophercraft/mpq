package cmd

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var hash_functions = map[string]func() hash.Hash{
	"md5":    md5.New,
	"sha1":   sha1.New,
	"sha256": sha256.New,
	"sha512": sha512.New,
}

func supported_hash_functions() (names []string) {
	names = make([]string, len(hash_functions))
	i := 0
	for k := range hash_functions {
		names[i] = k
		i++
	}
	sort.Strings(names)
	return
}

// list_set_command represents the list command
var list_set_command = &cobra.Command{
	Use:   "list <mpq file|json chain file>",
	Short: "Lists a MPQ set using a chain json file",
	Long:  `Load multiple MPQ archives into one set, listing all contents in a merged view`,
	Run: func(cmd *cobra.Command, args []string) {
		checksum, err := cmd.Flags().GetString("hash-algorithm")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		var hash_function func() hash.Hash
		if checksum != "" {
			var ok bool
			hash_function, ok = hash_functions[checksum]
			if !ok {
				fmt.Println("no hash function", checksum)
				os.Exit(1)
			}
		}

		count := 0
		set := open_set(cmd, args)

		set_list, err := set.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		var data []byte

		for set_list.Next() {
			count++
			path := set_list.Path()

			if hash_function != nil {
				file, err := set.Open(path)
				if err != nil {
					fmt.Println(path, err)
					break
				}
				data, err = io.ReadAll(file)
				if err != nil {
					fmt.Println("archived file read failed failed:", path, err)
					break
				}
				h := hash_function()
				h.Write(data)
				hash_sum := h.Sum(nil)
				hash_sum_string := hex.EncodeToString(hash_sum)
				fmt.Println(hash_sum_string, path)
				file.Close()
			} else {
				fmt.Println(path)
			}
		}

		set_list.Close()
		set.Close()
	},
}

func init() {
	add_set_commands(list_set_command)
	list_set_command.Flags().StringP("hash-algorithm", "a", "", fmt.Sprintf("the hash algorithm to use. set to an empty string to skip hashing. supported hash functions include %s", strings.Join(supported_hash_functions(), ", ")))
	root_cmd.AddCommand(list_set_command)
}
