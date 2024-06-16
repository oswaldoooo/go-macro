package token

import (
	"errors"
	"fmt"
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
type FuncType struct {
	Self    string
	Name    string
	Comment types.Comment
	Params  []Type
	Results []Type
	Body    string
}
type Value interface {
	Type
	Value() string
}
type FieldType struct {
	Name_    string
	Tag_     string
	Comment_ types.Comment
}
type BasicType struct {
	Name_ string
	Type_ string
}

func str2basic(src string, dst *Type) {
	bt := BasicType{}
	index := strings.IndexByte(src, ':')
	if index >= 0 {
		bt.Name_ = src[:index]
		bt.Type_ = src[index+1:]
	} else {
		bt.Type_ = src
	}
	*dst = bt
}
func (b BasicType) Name() string {
	return b.Name_
}

func (b BasicType) Type() string {
	return b.Type_
}

func (t *FieldType) Name() string {
	if index := strings.IndexByte(t.Name_, ':'); index > 0 {
		return t.Name_[:index]
	}

	return t.Name_
}

func (t *FieldType) Type() string {
	if index := strings.IndexByte(t.Name_, ':'); index > 0 {
		return t.Name_[index+1:]
	}
	return ""
}
func (f *FuncType) From(ft types.Func) {
	f.Self = ft.Self
	f.Name = ft.Name
	f.Comment = ft.Comment
	f.Params = make([]Type, len(ft.Params))
	utils.SliceConvert(ft.Params, f.Params, str2basic)
	f.Results = make([]Type, len(ft.Resp))
	utils.SliceConvert(ft.Resp, f.Results, str2basic)

}

type Struct struct {
	Name  string
	Field []FieldType
	//not use to unmarshale only for marshale
	Ident   string
	Tag     string
	Comment types.Comment
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
		dst.Name_ = src.Name
		dst.Tag_ = src.Tag
		dst.Comment_ = src.Comment
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

func (s Struct) MarshalText() (text []byte, err error) {
	var content string
	if len(s.Ident) > 0 {
		content = fmt.Sprintln(strings.Join(s.Comment, "\n")+"\ntype", s.Name, s.Ident)
		text = []byte(content)
		return
	}
	content = strings.Join(s.Comment, "\n") + "\ntype " + s.Name + " struct{"
	var fstr []byte
	if len(s.Field) > 0 {
		content += "\n"
	}
	for _, f := range s.Field {
		fstr, err = f.MarshalText()
		if err != nil {
			return
		}
		content += string(fstr)
	}
	content += "}"
	text = []byte(content)
	return
}
func (f FieldType) MarshalText() (text []byte, err error) {
	var content string = strings.ReplaceAll(f.Name_, ":", " ")
	if len(f.Comment_) > 0 {
		content = strings.Join(f.Comment_, "\n") + "\n" + content
	}
	if len(f.Tag_) > 0 {
		content += " " + f.Tag_
	}
	text = []byte(content + "\n")

	return

}
func (f FuncType) MarshalText() (text []byte, err error) {
	var content string = "func "
	if len(f.Comment) > 0 {
		content = strings.Join(f.Comment, "\n") + "\n" + content
	}
	if len(f.Self) > 0 {
		content += f.Self
	}
	content += f.Name + "("
	params := make([]string, len(f.Params))
	utils.SliceConvert(f.Params, params, func(src Type, dst *string) {
		*dst = src.Name() + " " + src.Type()
	})
	content += strings.Join(params, ",")
	content += ")"
	if len(f.Results) > 0 {
		params = make([]string, len(f.Results))
		utils.SliceConvert(f.Results, params, func(src Type, dst *string) {
			*dst = src.Name() + " " + src.Type()
		})
		content += "(" + strings.Join(params, ",") + ")"
	}
	content += "{"
	if len(f.Body) > 0 {
		content += "\n" + f.Body
	}
	content += "}\n"
	text = []byte(content)
	return
}

type PackageDecalre struct {
	PkgName string
	Import  []string
	Target_ string
}

func (p PackageDecalre) MarshalText() (text []byte, err error) {
	if len(p.PkgName) == 0 {
		err = errors.New("not set pkg name")
		return
	}
	content := "package " + p.PkgName + "\n"
	if len(p.Import) > 0 {
		content += "import(\n"
		for _, i := range p.Import {
			if len(i) == 0 {
				continue
			}
			index := strings.IndexByte(i, ':')
			if index >= 0 {
				content += i[:index] + " "
				i = i[index+1:]
			}
			content += "\"" + i + "\"\n"
		}
		content += ")\n"
	}
	text = []byte(content)
	return
}
func (p PackageDecalre) Target() string {
	if len(p.Target_) == 0 {
		panic("not set target go file")
	}
	return p.Target_
}
