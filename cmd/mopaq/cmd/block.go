package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Gophercraft/mpq"
	"github.com/spf13/cobra"
)

var block_cmd = &cobra.Command{
	Use:   "block [mpq file] [block index]",
	Short: "dump the indexed file block to stdout",
	Run:   run_block_cmd,
}

func init() {
	root_cmd.AddCommand(block_cmd)
}

func run_block_cmd(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		cmd.Help()
		os.Exit(1)
	}
	block_index, err := strconv.ParseUint(args[1], 10, 32)
	if err != nil {
		cmd.Help()
		os.Exit(1)
	}
	archive, err := mpq.Open(args[0])
	if err != nil {
		fmt.Fprintln(os.Stdout, "cannot open archive:", err)
		os.Exit(1)
	}

	block_entry, err := archive.BlockTableIndex(uint32(block_index))
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}

	archive.Close()

	block_position := archive.Position() + int64(block_entry.Position)
	block_size := block_entry.BlockSize

	archive_file, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}

	block_data := make([]byte, block_size)

	if _, err := archive_file.ReadAt(block_data, block_position); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}

	os.Stdout.Write(block_data)
}
