package token

import (
	"strings"

	"github.com/oswaldoooo/bgo/types"
	"github.com/oswaldoooo/go-macro/internal/utils"
)

type Ident struct {
	Name  string
	Value string
}

type Type interface {
	Name() string
	Type() string
}
type Value interface {
	Type
	Value() string
}
type FieldType struct {
	name    string
	tag     string
	comment types.Comment
}

func (t *FieldType) Name() string {
	if index := strings.IndexByte(t.name, ':'); index > 0 {
		return t.name[:index]
	}

	return t.name
}

func (t *FieldType) Type() string {
	if index := strings.IndexByte(t.name, ':'); index > 0 {
		return t.name[index+1:]
	}
	return ""
}

type Struct struct {
	Name  string
	Field []FieldType
}
type Const struct {
	name string
	val  string
}

func NewConst(n, v string) Const {
	return Const{name: n, val: v}
}
func NewVariable(n, v string) Variable {
	return Variable(NewConst(n, v))
}

type Variable Const

func (s *Struct) From(from types.Struct) {
	if len(from.Ident) > 0 {
		s.Name = from.Ident
		return
	}
	s.Name = from.Name
	s.Field = make([]FieldType, len(from.Fields))
	utils.SliceConvert(from.Fields, s.Field, func(src types.Field, dst *FieldType) {
		dst.name = src.Name
		dst.tag = src.Tag
		dst.comment = src.Comment
	})
}
func (t Const) Name() string {
	if index := strings.IndexByte(t.name, ':'); index > 0 {
		return t.name[:index]
	}

	return t.name
}

func (t Const) Type() string {
	if index := strings.IndexByte(t.name, ':'); index > 0 {
		return t.name[index+1:]
	}
	return ""
}

func (t Const) Value() string {
	return t.val
}
