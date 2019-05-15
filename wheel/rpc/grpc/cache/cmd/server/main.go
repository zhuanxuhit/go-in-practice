package main

import (
	"fmt"
	"github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/cache/pkg/cmd"
	"os"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
