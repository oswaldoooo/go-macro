package analyzer

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/oswaldoooo/bgo/types"
	"github.com/oswaldoooo/go-macro/token"
)

// analyze group
func (a *Analyzer) analyze_group(c context, group *types.Group[[]token.Value]) error {
	cmts := c.parseComments(group.Comments)
	for _, cmt := range cmts {
		if !cmt.IsValid() {
			continue
		}
		err := a.activeValue(c, cmt, group.Members)
		if err != nil {
			return err
		}
	}
	return nil
	// panic("not implement")
}
func (a *Analyzer) activeValue(c context, cmt Comment, g []token.Value) error {
	c.flag |= IgnoreErr
	c.val = cmt.Self
	if _, ok := a.repo[cmt.Self]; !ok {
		return errors.New("not found macro function " + cmt.Self)
	}
	f := a.repo[cmt.Self]
	groupParams, err := a.parseGroupParam(cmt.Params)
	if err != nil {
		return err
	}
	if groupParams.begin == -1 {
		groupParams.begin = 0
	}
	if groupParams.end == -1 {
		groupParams.end = len(g)
	}
	groupParams.others = append(groupParams.others, g[groupParams.begin:groupParams.end])
	// f.getFuncParamRequire()
	vlist := make([]reflect.Type, f.tp.NumIn())
	for i := range vlist {
		vlist[i] = f.tp.In(i)
	}
	vvlist, err := groupParams.try_into_params(a, c, vlist)
	if err != nil {
		return err
	}
	results := f.vl.Call(vvlist)
	for _, v := range results {
		if vstr, ok := v.Interface().(string); ok {
			a.appendToTail = append(a.appendToTail, vstr)
			continue
		}
		c.eprintf("for the time being, only direct appends are supported, and structured data is not returned")
	}
	return nil
}

type group_param struct {
	others []any
	//required params,-1 is unlimited
	begin int //group begin
	end   int //group end
}

func (a *Analyzer) parseGroupParam(params []string) (result group_param, err error) {
	var i, lasti int
	lasti = -1
	var (
		val any
	)
	for ; i < len(params); i++ {
		val = try_into(params[i])
		if params[i] == "*" {
			goto setval
		}
		if _, ok := val.(int64); ok {
			if lasti != -2 {
				goto setval
			}
		}
		result.others = append(result.others, val)
		continue
	setval:
		if lasti == -1 {
			lasti = i
			result.begin, err = topos(val)
			if err != nil {
				return
			}
		} else if i-lasti != 1 {
			err = errors.New("invalid params input. see examples/alias/example.go")
			return
		} else {
			result.end, err = topos(val)
			if err != nil {
				return
			}
			lasti = -2
		}
	}
	return
}
func topos(s any) (int, error) {
	if s == "*" {
		return -1, nil
	} else if v, ok := s.(int64); ok {
		return int(v), nil
	}
	return 0, errors.New(fmt.Sprint(s, "is not position. see examples/alias/example.go"))
}
func (g group_param) try_into_params(a *Analyzer, c context, target []reflect.Type) (results []reflect.Value, err error) {
	if len(g.others) != len(target) {
		err = errors.New("macro function " + c.val.(string) + " params error need " +
			strconv.Itoa(len(target)) + " provide " + strconv.Itoa(len(g.others)))
		return
	}
	results = make([]reflect.Value, len(target))
	for i := range results {
		val := reflect.ValueOf(g.others[i])
		if !val.Type().ConvertibleTo(target[i]) {
			err = errors.New("param can't not convertable")
			return
		}
		results[i] = val.Convert(target[i])
	}
	return
}
