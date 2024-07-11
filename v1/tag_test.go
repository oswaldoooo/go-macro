package v1_test

import (
	"encoding/json"
	"regexp"
	"testing"

	v1 "github.com/oswaldoooo/go-macro/v1"
)

// simple express
const (
	_tag_exp  = `(([\w\d]{1,}):"([\w\d\,\;]{0,})"){0,}`
	testwords = `json:"one,omitempty" yaml:"omitempty"`
)

var (
	_tag_reg = regexp.MustCompile(_tag_exp)
)

func TestReg(t *testing.T) {
	tg := *v1.ParseTag("`" + testwords + "`")
	content, _ := json.MarshalIndent(tg, "", "   ")
	t.Log(string(content))
	tg.Set("yaml", "jim")
	tg.Set("param", "jims")
	t.Log(tg.String())
}
