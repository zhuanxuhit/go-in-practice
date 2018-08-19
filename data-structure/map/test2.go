package main

import (
	"fmt"
	"unsafe"
)

func fn(m map[int]int) {
	fmt.Printf("%p -> %p : %d\n",&m, m, unsafe.Sizeof(m))
	m = make(map[int]int,100) // 分配内存
	fmt.Printf("%p -> %p : %d\n",&m, m, unsafe.Sizeof(m))
}

func main() {
	var m map[int]int
	fn(m)
	fmt.Println(m == nil)
}