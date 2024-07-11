package main

// call
// go-macro: override_struct(MyStruct)

type MyStruct struct {
	Name string `json:"Name"`
	Info string `json:"Info"`
	Kind int    `json:"Kind"`
}
type Enum uint8

func (s MyStruct) Getname() {}
func (s MyStruct) Close() error {
	return nil
}

const (
	Invalid Enum = iota + 1 //3
	Start                   //3
	Running                 //4
	End                     //5
)

type Stringer interface {
	String() string
}
