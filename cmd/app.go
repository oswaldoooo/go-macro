package main

import (
	"fmt"
	"os"

	"github.com/oswaldoooo/go-macro/analyzer"
	"github.com/oswaldoooo/go-macro/builder"
	"github.com/oswaldoooo/go-macro/token"
)

func throw(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error "+err.Error())
		os.Exit(-1)
	}
}
func main() {
	az, err := analyzer.NewAnalyzer("examples/example.go")
	throw(err)
	err = az.Analyze(analyzer.Normal)
	throw(err)
	err = builder.NewBuilder(az).Build(builder.Normal)
	throw(err)
}
func init() {
	analyzer.Register("decode", decode)
}

// example macro function
// decode
func decode(src token.Struct) (content string, out token.Struct) {
	out = src
	content = fmt.Sprintf("type Option%s %s", src.Name, src.Name)
	return
}
