package main

import (
	"flag"
	"fmt"
	log "github.com/golang/glog"
	"github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/v2/method"
	"net/rpc"
)

func main() {
	flag.Parse()
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call
	args := &method.Args{
		A: 7,
		B: 8,
	}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)

	//// Asynchronous call
	//quotient := new(Quotient)
	//divCall := client.Go("Arith.Divide", args, quotient, nil)
	//replyCall := <-divCall.Done	// will be equal to divCall
	//// check errors, print, etc.
}
