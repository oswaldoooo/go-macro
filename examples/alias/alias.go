package main

import (
	"fmt"
	"os"

	"github.com/oswaldoooo/bgo/types"
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
func init() {
	analyzer.Register("newobj", buildNewObject)
	analyzer.Register("printmethod", printMethod)
}
func printMethod(src []token.FuncType) {
	fmt.Println("find method", len(src))
	for _, f := range src {
		fmt.Println(f.Name, f.Params, f.Results)
	}
}
func buildNewObject(name string) (builder.Build[token.Struct], builder.Build[token.FuncType]) {
	return builder.Build[token.Struct]{
			Comment: types.Comment{"//this is a normal comment"},
			Data: token.Struct{
				Name: name,
				Field: []token.FieldType{
					token.FieldType{
						Name_:    "info:string",
						Tag_:     "`json:\"information\"`",
						Comment_: types.Comment{"//information"},
					},
				},
			},
		}, builder.Build[token.FuncType]{
			Data: token.FuncType{
				Self: "(self *" + name + ")",
				Name: "Getinfo",
				Results: []token.Type{
					token.BasicType{
						Type_: "string",
					},
				},
				Body: "return self.info",
			},
		}
}
