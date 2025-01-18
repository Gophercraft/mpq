package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Gophercraft/mpq"
	"github.com/spf13/cobra"
)

func open_set(cmd *cobra.Command, args []string) (set *mpq.Set) {
	if len(args) == 0 {
		cmd.Help()
		os.Exit(1)
	}

	unary_filepath := args[0]
	is_chain_json := false
	if filepath.Ext(unary_filepath) == ".json" {
		is_chain_json = true
	}

	if is_chain_json {
		chainjsondata, err := os.ReadFile(unary_filepath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		working_directory, err := cmd.Flags().GetString("working-directory")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		working_directory = filepath.Clean(working_directory)
		if err = os.Chdir(working_directory); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		var chain []string
		if err = json.Unmarshal(chainjsondata, &chain); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		set, err = mpq.GlobSet(chain...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		var err error
		set, err = mpq.GlobSet(unary_filepath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	return
}
