package analyzer

import (
	"errors"
	"reflect"
)

//see example.go.Self is decode Greeter is first Params

type Comment struct {
	Self   string
	Params []string
}
type funcinfo struct {
	tp reflect.Type
	vl reflect.Value
}

var (
	repo map[string]funcinfo = make(map[string]funcinfo)
)

func Register(name string, input any) error {
	tp := reflect.TypeOf(input)
	if tp.Kind() != reflect.Func {
		return errors.New("register must be func")
	}
	if _, ok := repo[name]; ok {
		return errors.New("duplicate error " + name)
	}
	repo[name] = funcinfo{
		tp: tp,
		vl: reflect.ValueOf(input),
	}
	return nil
}

func (c Comment) IsValid() bool {
	return len(c.Self) > 0
}

type paramRequire struct{}

func (p paramRequire) Match(in []reflect.Type) bool {

	return false
}

// get macro function param info
func (f funcinfo) getFuncParamRequire() (req paramRequire) {

	return
}

func GetMacroFuncNames() (result []string) {
	for i := range repo {
		result = append(result, i)
	}
	return
}
