package main

// call

// go-macro: printmethod(MyStruct?)
type MyStruct struct {
	Name string
	Info string
	Kind int
}
type Enum uint8

func (s MyStruct) Getname() {}
func (s MyStruct) Close() error {
	return nil
}

// go-macro: newobj(itest)
// go-macro: enum2str(*,*)
const (
	Invalid Enum = iota
	Start
	Running
	End
)
