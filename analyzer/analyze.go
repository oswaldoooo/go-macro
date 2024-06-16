package analyzer

import (
	"encoding"
	"errors"
	"go/format"
	stdparser "go/parser"
	"go/token"
	"os"
	"reflect"
	"strconv"
	"strings"
	_ "unsafe"

	"github.com/oswaldoooo/bgo/parser"
	"github.com/oswaldoooo/bgo/types"
	"github.com/oswaldoooo/go-macro/internal/utils"
	gtoken "github.com/oswaldoooo/go-macro/token"
)

type Analyzer struct {
	*types.Packages
	repo         map[string]funcinfo
	appendToTail []string
	rpath        string
}
type context struct {
	flag uint8
	pkgs *types.Package
	val  any
}

func NewAnalyzer(rpath string) (*Analyzer, error) {
	fs := token.NewFileSet()
	astf, err := stdparser.ParseFile(fs, rpath, nil, stdparser.AllErrors|stdparser.ParseComments)
	if err != nil {
		return nil, err
	}
	var result Analyzer
	result.rpath = rpath
	result.repo = repo
	result.Packages = types.NewPackages()
	err = parser.Parse(astf, result.Packages)
	return &result, err
}

// link the target go file and macro info
func (a *Analyzer) Analyze(_flag uint8) error {
	ctx := context{flag: _flag}
	for k, v := range a.Packages.Pkgs {
		ctx.pkgs = v
		err := a.analyze(k, ctx)
		if err != nil {
			return err
		}
	}
	//analyze group

	for _, g := range a.Packages.ConstGroup {
		vg := types.Group[[]gtoken.Value]{
			Comments: g.Comments,
			Members:  make([]gtoken.Value, len(g.Members)),
		}
		utils.SliceConvert(g.Members, vg.Members, func(src types.Const, dst *gtoken.Value) {
			vv := gtoken.NewConst(src.Name, src.Value)
			*dst = vv
		})
		err := a.analyze_group(ctx, &vg)
		if err != nil {
			return err
		}
	}
	return nil
}
func (a *Analyzer) analyze(_ string, ctx context) error {

	for _, t := range ctx.pkgs.Struct {
		cmts := ctx.parseComments(t.Comment)
		err := a.activeStructs(ctx, cmts)
		if err != nil {
			return err
		}
	}
	//todo: implement func,variable,const macro
	// for _, f := range ctx.pkgs.Func {

	// }
	return nil
}

// active the macro function to target
func (a *Analyzer) activeStructs(c context, cmts []Comment) error {
	for _, cmt := range cmts {
		if !cmt.IsValid() {
			continue
		}
		err := a.active_struct(c, cmt.Self, cmt.Params)
		if err != nil {
			return err
		}
	}
	return nil
}
func (a *Analyzer) active_struct(c context, self string, params []string) error {
	// fmt.Println("get macro func " + self)
	f, ok := a.repo[self]
	if !ok {
		return errors.New("not found macro function " + self)
	}
	if f.tp.NumIn() != len(params) {
		return errors.New("macro function " + self + " params error need " +
			strconv.Itoa(f.tp.NumIn()) + " provide " + strconv.Itoa(len(params)))
	}
	//transfer params
	var fin = make([]reflect.Value, len(params))
	for i := range fin {
		if strings.HasSuffix(params[i], "?") {
			fin[i] = a.matchFuncs(c, params[i][:len(params[i])-1])
			continue
		}
		v := try_into(params[i])
		if iv, ok := v.(int64); ok {
			fin[i] = c.convert(iv, f.tp.In(i))
		} else if fv, ok := v.(float64); ok {
			fin[i] = c.convert(fv, f.tp.In(i))
		} else if bv, ok := v.(bool); ok {
			fin[i] = c.convert(bv, f.tp.In(i))
		} else {
			fin[i] = a.tryGetType(c, v.(string), f.tp.In(i))
		}
		if !fin[i].IsValid() {
			return errors.New("param " + params[i] + " type error")
		}
	}
	results := f.vl.Call(fin)
	//analyze result
	for _, v := range results {
		if vstr, ok := v.Interface().(string); ok {
			a.appendToTail = append(a.appendToTail, vstr)
			continue
		} else if marshaler, ok := v.Interface().(encoding.TextMarshaler); ok {
			content, err := marshaler.MarshalText()
			if err != nil {
				return err
			}
			newcontent, err := format.Source(content)
			if err == nil {
				content = newcontent
			}
			a.appendToTail = append(a.appendToTail, string(content))
			continue
		}
	}
	return nil
}

// if not found return self
func (a *Analyzer) tryGetType(c context, name string, tp reflect.Type) (result reflect.Value) {
	if st, ok := c.pkgs.Struct[name]; ok {
		var rr gtoken.Struct
		rr.From(st)
		result = reflect.ValueOf(rr)
		if !result.Type().ConvertibleTo(tp) {
			result = reflect.Value{}
		} else {
			result = result.Convert(tp)
		}
		return
	}
	return reflect.ValueOf(name)
	// panic("not implement other type")
}
func (a *Analyzer) matchFuncs(c context, name string) (result reflect.Value) {
	var ans []gtoken.FuncType
	var ff gtoken.FuncType
	for _, f := range c.pkgs.Func {
		if len(f.Self) == 0 {
			continue
		}
		if index := strings.IndexByte(f.Self, ' '); index >= 0 {
			f.Self = f.Self[index+1:]
		}
		if f.Self != name {
			continue
		}
		ff.From(f)
		ans = append(ans, ff)
	}
	result = reflect.ValueOf(ans)
	return
}

//go:linkname build_additional github.com/oswaldoooo/go-macro/builder.analyze_build
func build_additional(a *Analyzer) error {
	// fmt.Println("append to tail ", len(a.appendToTail))
	if len(a.appendToTail) == 0 {
		return nil
	}
	result := strings.Join(a.appendToTail, "\n\n")
	f, err := os.OpenFile(a.rpath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(result)
	return err
}
