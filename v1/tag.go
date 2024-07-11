package v1

import (
	"regexp"
	"strings"
)

type tagfield struct {
	name  string
	value string
}
type Tag []tagfield

// simple express
const (
	_tag_exp = `(([\w\d]{1,}):"([\w\d\,\;]{0,})"){0,}`
)

var (
	_tag_reg = regexp.MustCompile(_tag_exp)
)

func ParseTag(s string) *Tag {
	var t Tag
	if len(s) > 0 {
		s = s[1 : len(s)-1]
	}
	for _, v := range _tag_reg.FindAllStringSubmatch(s, -1) {
		if len(v) > 0 && len(v[2]) > 0 {
			t = append(t, tagfield{name: v[2], value: v[3]})
		}
	}
	return &t
}
func (t *Tag) Get(key string, defaul string) string {
	for i := range *t {
		if (*t)[i].name == key {
			return (*t)[i].value
		}
	}
	return defaul
}
func (t *Tag) Set(k, v string) {
	for i := range *t {
		if (*t)[i].name == k {
			(*t)[i].value = v
			return
		}
	}
	(*t) = append((*t), tagfield{name: k, value: v})
}
func (t Tag) String() string {
	var (
		ans []string = make([]string, len(t))
		i   int
	)
	for k := range t {
		ans[i] = t[k].name + `:"` + t[k].value + `"`
		i++
	}
	if len(ans) > 0 {
		return "`" + strings.Join(ans, " ") + "`"
	}
	return ""
}
