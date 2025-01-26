package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var read_cmd = &cobra.Command{
	Use:   "read [mpq files/chain json] [archived file to be read]",
	Short: "read a File from the Set, copying its contents to stdout",
	Run:   run_read_cmd,
}

func init() {
	root_cmd.AddCommand(read_cmd)
}

func run_read_cmd(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Help()
		os.Exit(1)
	}

	set := open_set(cmd, args[:len(args)-1])
	file, err := set.Open(args[len(args)-1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err = io.Copy(os.Stdout, file); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file.Close()
}
