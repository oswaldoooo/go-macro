package examples_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/oswaldoooo/bgo"
)

func TestShow(t *testing.T) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "example.go", nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	ast.Print(fs, f)
}

func TestPrint(t *testing.T) {
	r, p, err := bgo.Parse("example.go")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile("example.json", bgo.Debug(*p), 0644)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile("example.ast", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = ast.Fprint(f, r.Fs, r.Astree, nil)
	if err != nil {
		t.Fatal(err)
	}
}
