package main

import (
	"fmt"
	"os"

	"github.com/Gophercraft/mpq"
)

func main() {
	path := os.Args[1]

	archive, err := mpq.Open(path)
	if err != nil {
		panic(err)
	}

	list, err := archive.List()
	if err != nil {
		panic(err)
	}

	for list.Next() {
		fmt.Println("List => ", list.Path())
	}
	list.Close()
	archive.Close()
}
