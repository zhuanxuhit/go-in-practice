package main

import (
	"unsafe"
)

func main() {
	s := "abcdefg"
	// 16 7
	println(unsafe.Sizeof(s), len(s))
	// 7 是字符串长度，为什么 sizeof是 16，因为是一个指针加长度，都是8字节，2个8字节就是16了
}
