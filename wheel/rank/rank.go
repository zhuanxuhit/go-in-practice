package main

import (
	"github.com/XanthusL/zset"
	"log"
)

func main() {
	set := zset.New()
	set.Set(1.0, 1, "test1")
	set.Set(2.0, 2, "test2")
	userRank, _, _ := set.GetRank(2, false)
	log.Printf("%+v", userRank)
}
