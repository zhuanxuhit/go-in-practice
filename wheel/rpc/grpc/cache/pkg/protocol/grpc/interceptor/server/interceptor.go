package server

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

func WithServerTimeInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(serverTimeInterceptor)
}

func serverTimeInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("invoke server method=%s duration=%s error=%v", info.FullMethod,
		time.Since(start), err)
	return resp, err
}
