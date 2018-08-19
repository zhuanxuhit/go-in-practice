package main

import (
	"unsafe"
)

func toString(b []byte) string {
	return string(b)
}

func unsafeToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))//b指针的数据结构
}

func main() {
	b := []byte("hello, world!")
	s1 := toString(b)
	s2 := unsafeToString(b)

	println(s1 == s2)
	println(s1, s2)
}