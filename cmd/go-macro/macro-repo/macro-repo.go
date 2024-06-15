package macrorepo

import (
	"fmt"
	"go/format"

	"github.com/oswaldoooo/go-macro/analyzer"
	"github.com/oswaldoooo/go-macro/token"
)

func Init() {
	analyzer.Register("enum2str", enum2str)
}
func enum2str(name string, src []token.Value) (content string) {
	// srclist := make([]string, len(src))
	for _, v := range src {
		content += "    " + v.Name() + ":\"" + v.Name() + "\",\n"
	}
	content = fmt.Sprintf("var %s =map[%s]string{\n%s\n}", name, src[0].Type(), content)
	con, err := format.Source([]byte(content))
	if err != nil {
		fmt.Println("warn format failed " + err.Error())
	} else {
		content = string(con)
	}
	return
}
