package main

import (
	"strings"
	"git.apache.org/thrift.git/lib/go/thrift"
	"log"
	"fmt"
	"go-in-practice/thrift/gen-go/echo"
	"context"
)

//内部原理
type FormatDataImpl struct {}

func (fdi *FormatDataImpl) DoFormat(ctx context.Context, data *echo.Data) (r *echo.Data, err error){
	var rData echo.Data
	rData.Text = strings.ToUpper(data.Text)

	return &rData, nil
}

const (
	HOST = "localhost"
	PORT = "8080"
)

func main() {

	handler := &FormatDataImpl{}
	processor := echo.NewFormatDataProcessor(handler)
	serverTransport, err := thrift.NewTServerSocket(HOST + ":" + PORT)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	fmt.Println("Running at:", HOST + ":" + PORT)
	server.Serve()
}

