package builder

import (
	_ "unsafe"

	"github.com/oswaldoooo/go-macro/analyzer"
)

type Builder struct {
	a *analyzer.Analyzer
}

// flag
const (
	Normal uint8 = 1 << iota
	NoWarn
	IgnoreErr
)

func NewBuilder(a *analyzer.Analyzer) *Builder {
	b := Builder{
		a: a,
	}
	return &b
}

// build with flag
func (b *Builder) Build(_flag uint8) error {
	return analyze_build(b.a)
}

//go:linkname analyze_build
//go:noescape
func analyze_build(*analyzer.Analyzer) error
