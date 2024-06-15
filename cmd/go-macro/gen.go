package main

import (
	"errors"
	"path"

	"github.com/oswaldoooo/go-macro/analyzer"
	"github.com/oswaldoooo/go-macro/builder"
	"github.com/spf13/cobra"
)

func newGenCmd() *cobra.Command {
	var cmd = cobra.Command{
		Use:          "gen",
		Short:        "generate golang code by macro call",
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				throw(errors.New("not set target go file"))
			}
			nowarn, _ := cmd.Flags().GetBool("nowarn")
			inputdir, _ := cmd.Flags().GetString("input")
			_flag := analyzer.Normal
			if nowarn {
				_flag |= analyzer.NoWarn
			}
			pdir := pather(inputdir)
			for _, t := range args {
				a, err := analyzer.NewAnalyzer(pdir.Get(t))
				throw(err)
				throw(a.Analyze(_flag))
				throw(builder.NewBuilder(a).Build(_flag))
			}
		},
	}
	cmd.Flags().StringP("input", "i", "", "input directory")
	cmd.Flags().BoolP("nowarn", "w", false, "will not show warning msg")
	return &cmd
}

type pather string

func (p pather) Get(s string) string {
	if len(p) == 0 {
		return s
	}
	return path.Join(string(p), s)
}
