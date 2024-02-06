package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Gophercraft/mpq"
	"github.com/spf13/cobra"
)

// listSetCmd represents the listSet command
var listSetCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists a MPQ set",
	Long:  `Load multiple MPQ archives into one set, listing all contents in a merged view`,
	Run: func(cmd *cobra.Command, args []string) {
		working_directory, err := cmd.Flags().GetString("working-directory")
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
		set_list, err := set.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		count := 0

		for set_list.Next() {
			count++
			path := set_list.Path()
			fmt.Println("List =>", path)
		}

		set_list.Close()
		set.Close()
		fmt.Println(count, "files listed!")
	},
}

func init() {
	rootCmd.AddCommand(listSetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listSetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listSetCmd.Flags().StringP("chain-json", "c", "", "load a list of MPQ globs from a JSON file")
	listSetCmd.Flags().StringP("working-directory", "w", "", "working directory")
}
