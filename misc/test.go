package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	a := r.Perm(200)[:5]
	fmt.Println(a)
}
