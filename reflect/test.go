package main

import (
	"fmt"
	"reflect"
)

type Test struct {
	A int
	B string
	C float32
}

func main() {
	test := new(Test)
	fmt.Println(reflect.TypeOf(test))
	fmt.Println(reflect.ValueOf(test))
	fmt.Println(reflect.Indirect(reflect.ValueOf(test)).Type().Name())
}
