package main

import (
	"fmt"
	"os"

	"github.com/joseluisq/enve/cmd"
)

func main() {
	if err := cmd.Execute(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
