package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gophercraft/mpq"
	"github.com/spf13/cobra"
)

var export_cmd = &cobra.Command{
	Use:   "export",
	Short: "exports an MPQ set to a directory",
	Long:  `Load multiple MPQ archives into one set, then exports all of its contents to `,
	Run: func(cmd *cobra.Command, args []string) {
		cd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		// get working directory
		working_directory, err := cmd.Flags().GetString("working-directory")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		// get export directory
		export_directory, err := cmd.Flags().GetString("export-directory")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
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
		os.Chdir(working_directory)
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
		os.Chdir(cd)

		set_list, err := set.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		base_export_path := export_directory

		for set_list.Next() {
			count++
			path := set_list.Path()
			output_path_parts := strings.Split(path, "\\")
			output_path := filepath.Join(append([]string{base_export_path}, output_path_parts...)...)
			file, err := set.Open(path)
			if err != nil {
				fmt.Println(path, err)
				break
			}
			dir_containing_output := filepath.Dir(output_path)
			if err = os.MkdirAll(dir_containing_output, 0700); err != nil {
				fmt.Println(path, err)
				break
			}
			fmt.Println(output_path)
			output_file, err := os.OpenFile(output_path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0700)
			if err != nil {
				fmt.Println(path, err)
				break
			}
			if _, err = io.Copy(output_file, file); err != nil {
				fmt.Println(path, err)
				break
			}
			file.Close()
		}

		set_list.Close()
		set.Close()
	},
}

func init() {
	root_cmd.AddCommand(export_cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listSetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	export_cmd.Flags().StringP("chain-json", "c", "", "load a list of MPQ globs from a JSON file")
	export_cmd.Flags().StringP("working-directory", "w", "", "working directory")
	export_cmd.Flags().StringP("export-directory", "o", "", "output directory")
}
