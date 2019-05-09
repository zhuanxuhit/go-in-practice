package main

import (
	"context"
	"flag"
	pb "github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/v4/helloworld"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

var (
	name string
)

func init() {
	flag.StringVar(&name, "name", defaultName, "name")
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	// Contact the server and print out its response
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
