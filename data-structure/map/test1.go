package main

import "unsafe"

func main() {
	max := 1000
	m := make(map[int]int,1000)
	println(unsafe.Sizeof(m),len(m))
	for i := 0; i < max; i++ {
		m[i] = i
	}
}
