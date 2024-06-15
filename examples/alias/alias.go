package main

import (
	"fmt"
	"os"

	"github.com/oswaldoooo/go-macro/analyzer"
	"github.com/oswaldoooo/go-macro/builder"
	"github.com/oswaldoooo/go-macro/token"
)

func main() {
	az, err := analyzer.NewAnalyzer("example.go")
	throw(err)
	err = az.Analyze(analyzer.Normal)
	throw(err)
	err = builder.NewBuilder(az).Build(builder.Normal)
	throw(err)
}

func alias(aliasname string, src token.Struct) (content string) {
	content = fmt.Sprintf("type %s %s", aliasname, src.Name)
	return
}

func enum2str(enum []token.Value) (result string) {
	result = ""
	for _, v := range enum {
		fmt.Println("const ", v.Name(), v.Type(), v.Value())
	}
	return
}
func init() {
	analyzer.Register("alias", alias)
	analyzer.Register("enum2str", enum2str)
}
func throw(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error "+err.Error())
		os.Exit(-1)
	}
}
