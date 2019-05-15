package main

import (
	"context"
	"flag"
	"github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/cache/pkg/api/v1"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	// get configuration
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := v1.NewCacheClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := "v1"
	val := "我是v1版本store"
	_, err = c.Store(ctx, &v1.StoreReq{
		Key: key,
		Val: []byte(val),
	})
	if err != nil {
		log.Fatalf("Store failed: %v", err)
	}
	resp, err := c.Get(ctx, &v1.GetReq{Key: key})
	if err != nil {
		log.Fatalf("Store failed: %v", err)
	}
	if string(resp.Val) != val {
		log.Fatalf("wanted: %s, but get:%s", val, string(resp.Val))
	}
}
