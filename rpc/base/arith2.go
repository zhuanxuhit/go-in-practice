package main

import (
	"net/rpc"
	"net"
	"log"
	"fmt"
	"awesomeProject/go-in-practice/rpc/proto"
)

func init() {
	serv := rpc.NewServer()
	arith := new(proto.Arith)
	// 注册了service，service名字就不指定就是用 arith
	serv.RegisterName("Arith",arith)

	lis, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatalln("listen error.")
	}
	// accept每个连接，然后在携程中处理请求
	go serv.Accept(lis)
}
func main() {
	// 客户端
	client, err := rpc.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	args := &proto.Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
}
