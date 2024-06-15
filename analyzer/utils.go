package analyzer

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/oswaldoooo/bgo/types"
)

// flag
const (
	Normal uint8 = 1 << iota
	NoWarn
	IgnoreErr
)
const Mark = "go-macro:"

// return 0 if not match
func (c context) parseComment(s string) (ans Comment) {
	s = s[2:]
	s = strings.TrimSpace(s)
	// fmt.Println("parse comment:", s)
	if strings.HasPrefix(s, Mark) {
		s = s[len(Mark):]
		// fmt.Println("parse go-macro:", s)
		s = strings.TrimSpace(s)
		index := strings.IndexByte(s, '(')
		if index <= 0 {
			c.eprintf("warning parse comment " + s + " failed. the macro format is like macro_name(param1,param2)")
			return
		}
		selfname := s[:index]
		s = s[index+1:]
		index = strings.IndexByte(s, ')')
		if index <= 0 {
			c.eprintf("warning parse comment " + s + " failed. the macro format is like macro_name(param1,param2)")
			return
		}
		params := strings.Split(s[:index], ",")
		for _, e := range params {
			if len(e) == 0 {
				c.eprintf(s[:index] + " parma " + e + " is null")
				return
			}
		}
		ans.Self = selfname
		ans.Params = params
	}
	return
}

func (c context) parseComments(s types.Comment) (ans []Comment) {
	ans = make([]Comment, len(s))
	for i := range s {
		ans[i] = c.parseComment(s[i])
	}
	return
}

func (c context) convert(input any, tp reflect.Type) reflect.Value {
	v := reflect.ValueOf(input)
	if !v.Type().ConvertibleTo(tp) {
		return reflect.Value{}
	}
	return v.Convert(tp)
}
func (c context) eprintf(format string, a ...any) {
	if !strings.HasPrefix(format, "\n") {
		format += "\n"
	}
	if c.flag&IgnoreErr == 0 {
		fmt.Fprintf(os.Stderr, format, a...)
		os.Exit(-1)
	}
	if c.flag&NoWarn == 0 {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}

func try_into(s string) any {
	val, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return val
	}
	fval, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return fval
	}
	bval, err := strconv.ParseBool(s)
	if err == nil {
		return bval
	}
	return s
}
