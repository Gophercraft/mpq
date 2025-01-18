package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var root_cmd = &cobra.Command{
	Use:   "mopaq",
	Short: "A MoPaQ extraction utility written in Go",
	Long:  "https://github.com/Gophercraft/mpq",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := root_cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	root_cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
