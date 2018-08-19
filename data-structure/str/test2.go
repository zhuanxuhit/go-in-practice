package main

import (
	"strings"
	"unsafe"
)

func main() {
	s := strings.Repeat("a",3)
	println(unsafe.Sizeof(s), len(s))
}
