# 1th Go-Macro

## What is Go-Macro?
**go-macro** is a macro diy package. You can define yourself special macro function to generate code.

## Example
### 1. Alias macro function like below
```go
import (
  "github.com/oswaldoooo/go-macro/token"
)
func alias(aliasname string, src token.Struct) (content string) {
	content = fmt.Sprintf("type %s %s", aliasname, src.Name)
	return
}

// use example see examples/alias/example.go
```
### 2.Enum take
```go
import (
  "github.com/oswaldoooo/go-macro/token"
)
func enum2str(enum []token.Value) (result string) {
	result = ""
	for _, v := range enum {
		fmt.Println("const ", v.Name(), v.Type(), v.Value())
	}
	return
}
```

### Notice
**1th Beta Edition is unstable now.**