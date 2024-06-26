package builder

import (
	"encoding"
	"strings"
	_ "unsafe"

	"github.com/oswaldoooo/bgo/types"
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

type Build[T encoding.TextMarshaler] struct {
	Target_ string        //if target is not null,it will be write to this file
	Comment types.Comment //comment
	Data    T
	Slice   []T
}

func (b Build[T]) Target() string {
	return b.Target_
}
func (b Build[T]) MarshalText() (text []byte, err error) {
	var content string
	if len(b.Comment) > 0 {
		content = strings.Join(b.Comment, "\n")
	}
	if len(b.Slice) > 0 {
		var totaltext []byte
		for _, ss := range b.Slice {
			text, err = ss.MarshalText()
			if err != nil {
				return
			}
			totaltext = append(totaltext, text...)
		}
		text = totaltext
	} else {
		text, err = b.Data.MarshalText()
		if err != nil {
			return
		}
	}

	content += string(text)
	text = []byte(content)
	return
}
