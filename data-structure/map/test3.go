package main

import (
	"fmt"
	"encoding/json"
)

func main() {
	var mm map[string]string
	aa := mm["aa"]
	fmt.Println("value",aa)
	return

	m := make([]*int, 0, 6)
	s, _ := json.Marshal(m)
	fmt.Println(string(s))
	a := 1
	for i := 0; i < 6; i ++ {
		m[i] = &a
	}
	fmt.Println(json.Marshal(m))
}
