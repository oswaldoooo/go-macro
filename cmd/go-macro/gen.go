package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path"
	"plugin"

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
			loadplugins()
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

type Plugins struct {
	XMLName xml.Name `xml:"plugins"`
	Plugins []Plugin `xml:"plugin"`
	Include []string `xml:"include"`
}
type Plugin struct {
	Path string `xml:"path,attr"`
	Name string `xml:",chardata"`
}

func loadplugins() {
	loadpluginconfig(path.Join(os.Getenv("HOME"), ".config", "go-macro", "plugins.xml"))
}

func loadpluginconfig(cpath string) {
	homepath := os.Getenv("HOME")
	content, err := os.ReadFile(cpath)
	if err != nil {
		return
	}
	var pl Plugins
	err = xml.Unmarshal(content, &pl)
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning plugin config file format error "+err.Error())
		return
	}
	homepath = path.Join(homepath, ".config", "go-macro", "plugins")
	for _, pp := range pl.Plugins {
		defaulthome := homepath
		if len(pp.Path) > 0 {
			defaulthome = pp.Path
		} else if len(pp.Name) == 0 {
			fmt.Fprintln(os.Stderr, "not set plugin name")
			continue
		}
		err = loadplugin(path.Join(defaulthome, pp.Name+".a"))
		if err != nil {
			fmt.Fprintln(os.Stderr, "load plugin "+pp.Name+" error "+err.Error())
		}
	}
	for _, inc := range pl.Include {
		loadpluginconfig(inc)
	}
}
func loadplugin(plugipath string) error {
	pl, err := plugin.Open(plugipath)
	if err != nil {
		return err
	}
	sp, err := pl.Lookup("MacroList")
	if err != nil {
		return err
	}
	macrolist, ok := sp.(*[]string)
	if !ok {
		return errors.New("macro list is not []string")
	}
	for _, m := range *macrolist {
		fs, err := pl.Lookup(m)
		if err != nil {
			return errors.New("not found " + m + " error " + err.Error())
		}
		err = analyzer.Register(m, fs)
		if err != nil {
			return errors.New("register macro function " + m + " error " + err.Error())
		}
	}
	return nil
}
