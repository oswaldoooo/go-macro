package main

// call

//go macro

type MyStruct struct {
	Name string
	Info string
	Kind int
}
type Enum uint8

// go-macro: enum2str(EnumStr,*,*)
const (
	Invalid Enum = iota
	Start
	Running
	End
)

var EnumStr = map[Enum]string{
	Invalid: "Invalid",
	Start:   "Start",
	Running: "Running",
	End:     "End",
}
