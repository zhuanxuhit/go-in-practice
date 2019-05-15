package grpc

import (
	"context"
	"github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/cache/pkg/api/v1"
	"github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/cache/pkg/protocol/grpc/interceptor/server"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func RunServer(ctx context.Context, v1API v1.CacheServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// register service
	s := grpc.NewServer(server.WithServerTimeInterceptor())
	v1.RegisterCacheServer(s, v1API)
	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC s...")
			s.GracefulStop()
			//<-ctx.Done()
		}
	}()

	// start gRPC s
	log.Println("starting gRPC s...")
	return s.Serve(listen)
}
