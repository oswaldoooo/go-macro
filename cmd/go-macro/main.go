package main

import (
	"fmt"
	"os"
	"runtime"

	macrorepo "github.com/oswaldoooo/go-macro/cmd/go-macro/macro-repo"
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Short: "go-macro 1.0 beta    go version: " + runtime.Version(),
}

func main() {
	macrorepo.Init()
	rootCmd.AddCommand(newGenCmd(), newListCommand())
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func throw(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error "+err.Error())
		os.Exit(-1)
	}
}
