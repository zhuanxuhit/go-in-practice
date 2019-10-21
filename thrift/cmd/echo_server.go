package main

import (
	"awesomeProject/go-in-practice/thrift/api/echo_thrift/gen-go/echo"
	"fmt"
	"context"
)

type EchoServerImp struct {
}

func (e *EchoServerImp) Echo(ctx context.Context, req *echo.EchoReq) (*echo.EchoRes, error) {
	fmt.Printf("message from client: %v\n", req.GetName())

	res := &echo.EchoRes{
		Msg: "hello " + req.GetName(),
	}

	return res, nil
}

func main() {
	//transport, err := thrift.NewTServerSocket(":3000")
	//if err != nil {
	//	panic(err)
	//}
	//
	//processor := echo.NewEchoProcessor(&EchoServerImp{})
	//server := thrift.NewTSimpleServer4(
	//	processor,
	//	transport,
	//	thrift.NewTBufferedTransportFactory(8192),
	//	thrift.NewTCompactProtocolFactory(),
	//)
	//if err := server.Serve(); err != nil {
	//	panic(err)
	//}
}
