package main

import (
	"fmt"
	"strings"

	"github.com/oswaldoooo/go-macro/analyzer"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {

	var cmd = cobra.Command{
		Use:   "list",
		Short: "show available resource list",
		Run: func(cmd *cobra.Command, args []string) {
			loadplugins()
			if ok, _ := cmd.Flags().GetBool("macro"); ok {
				fmt.Println(strings.Join(analyzer.GetMacroFuncNames(), "\n"))
			}
		},
	}
	cmd.Flags().Bool("macro", false, "show registered macro resource")
	return &cmd
}
