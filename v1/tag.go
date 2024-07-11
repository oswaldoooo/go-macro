package v1

import (
	"regexp"
	"strings"
)

type Tag map[string]string

// simple express
const (
	_tag_exp = `(([\w\d]{1,}):"([\w\d\,\;]{0,})"){0,}`
)

var (
	_tag_reg = regexp.MustCompile(_tag_exp)
)

func ParseTag(s string) *Tag {
	var t = make(Tag)
	if len(s) > 0 {
		s = s[1 : len(s)-1]
	}
	for _, v := range _tag_reg.FindAllStringSubmatch(s, -1) {
		t[v[2]] = v[3]
	}
	return &t
}
func (t *Tag) Get(key string, defaul string) string {
	if v, ok := (*t)[key]; ok {
		return v
	}
	return defaul
}
func (t *Tag) Set(k, v string) {
	(*t)[k] = v
}
func (t Tag) String() string {
	var (
		ans []string = make([]string, len(t))
		i   int
	)
	for k, v := range t {
		ans[i] = k + `:"` + v + `"`
		i++
	}
	if len(ans) > 0 {
		return "`" + strings.Join(ans, " ") + "`"
	}
	return ""
}
