package main

import (
	"github.com/oswaldoooo/bgo/types"
	"github.com/oswaldoooo/go-macro/analyzer"
	"github.com/oswaldoooo/go-macro/builder"
	"github.com/oswaldoooo/go-macro/token"
)

func init() {
	analyzer.Register("newobj", buildNewObject)
}
func buildNewObject(name string) builder.Build[token.Struct] {
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
	}
}
