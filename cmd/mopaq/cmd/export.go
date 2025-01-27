package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
		// get export directory
		export_directory, err := cmd.Flags().GetString("export-directory")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		set := open_set(cmd, args)
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

			file, err := set.Open(path)
			if err != nil {
				fmt.Println(path, err)
				break
			}
			file_data, err := io.ReadAll(file)
			if err != nil {
				fmt.Println(path, err)
				break
			}
			file.Close()
			file_data_hash := sha256.Sum256(file_data[:])
			write_new_file := true

			path = strings.ToLower(path)
			output_path_parts := strings.Split(path, "\\")
			output_path := filepath.Join(append([]string{base_export_path}, output_path_parts...)...)

			existing_file_data, err := os.ReadFile(output_path)
			if err == nil {
				existing_file_data_hash := sha256.Sum256(existing_file_data[:])
				if existing_file_data_hash == file_data_hash {
					write_new_file = false
				}
			}

			if write_new_file {
				fmt.Println("write:", output_path)
				dir_containing_output := filepath.Dir(output_path)
				if err = os.MkdirAll(dir_containing_output, 0700); err != nil {
					fmt.Println(path, err)
					break
				}
				if err = os.WriteFile(output_path, file_data, 0700); err != nil {
					fmt.Println(path, err)
					break
				}
			} else {
				fmt.Println("exists:", output_path)
			}
		}

		set_list.Close()
		set.Close()
	},
}

func init() {
	add_set_commands(export_cmd)
	export_cmd.Flags().StringP("export-directory", "o", "", "output directory")
	root_cmd.AddCommand(export_cmd)

}
