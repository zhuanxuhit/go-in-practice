package client

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)
func WithClientTimeInterceptor() grpc.DialOption {
	return grpc.WithUnaryInterceptor(timeInterceptor)
}

func timeInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...) // <==
	log.Printf("invoke remote method=%s duration=%s error=%v", method,
		time.Since(start), err)
	return err
}
