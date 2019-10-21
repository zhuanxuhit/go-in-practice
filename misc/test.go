package main

import (
	"fmt"
	"strings"
)

func main() {
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//a := r.Perm(200)[:5]
	//fmt.Println(a)
	//ctx := context.Background()
	//<-ctx.Done()
	p := "ds=ab,af=12,sd=4"
	spec := strings.SplitN(p, "=", 2)
	fmt.Printf("%+v",spec)
}
