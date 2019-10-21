package main

import (
	"fmt"
	"unsafe"
)

type S interface {
	Get() int
	Set(int)
}

type sInt struct {
	value int
}

func (s *sInt) Get() int {
	return s.value
}
func (s *sInt) Set(v int) {
	s.value = v
}

func main() {
	var s S
	s = &sInt{
		value: 10,
	}
	fmt.Println(s.Get())
	fmt.Println(unsafe.Sizeof(s))
}
