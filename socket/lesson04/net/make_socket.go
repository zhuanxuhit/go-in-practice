package main

import (
	"code.byted.org/gopkg/pkg/log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:12345")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for  {
		;
	}
}
