package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Gophercraft/mpq"
	"github.com/spf13/cobra"
)

// list_set_command represents the list command
var list_set_command = &cobra.Command{
	Use:   "list",
	Short: "Lists a MPQ set",
	Long:  `Load multiple MPQ archives into one set, listing all contents in a merged view`,
	Run: func(cmd *cobra.Command, args []string) {
		working_directory, err := cmd.Flags().GetString("working-directory")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		working_directory = filepath.Clean(working_directory)
		chainjson, err := cmd.Flags().GetString("chain-json")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		if chainjson == "" {
			fmt.Println("you need to supply a MPQ chain")
			os.Exit(0)
		}
		chainjsondata, err := os.ReadFile(chainjson)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		if err = os.Chdir(working_directory); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		var chain []string
		if err = json.Unmarshal(chainjsondata, &chain); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		set, err := mpq.GlobSet(chain...)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		count := 0

		set_list, err := set.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		var data []byte

		for set_list.Next() {
			count++
			path := set_list.Path()
			file, err := set.Open(path)
			if err != nil {
				fmt.Println(path, err)
				break
			}
			data, err = io.ReadAll(file)
			if err != nil {
				fmt.Println("readall failed:", path, err)
				break
			}
			sum256 := sha256.Sum256(data)
			shasum := hex.EncodeToString(sum256[:])
			fmt.Println(shasum, path)

			file.Close()
		}

		set_list.Close()
		set.Close()
	},
}

func init() {
	root_cmd.AddCommand(list_set_command)

	list_set_command.Flags().StringP("chain-json", "c", "", "load a list of MPQ globs from a JSON file")
	list_set_command.Flags().StringP("working-directory", "w", "", "working directory")
}
