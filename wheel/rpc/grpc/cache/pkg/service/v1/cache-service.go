package v1

import (
	"context"
	"github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/cache/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"sync"
)

type cacheServiceServer struct {
	store map[string][]byte
	mutex sync.Mutex
}

func NewCacheServiceServer() v1.CacheServer {
	server := new(cacheServiceServer)
	server.store = make(map[string][]byte)
	return server
}

func (server *cacheServiceServer) Store(ctx context.Context, req *v1.StoreReq) (*v1.StoreResp, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok1 := md["test-run"]; ok1 {
			if len(val) > 0 && val[0] == "1" {
				// 不进行请求，是一个测试，直接返回
				log.Print("test-run")
				return &v1.StoreResp{}, nil
			}
		}
	}
	server.mutex.Lock()
	server.store[req.Key] = req.Val
	server.mutex.Unlock()
	return &v1.StoreResp{}, nil
}

func (server *cacheServiceServer) Get(ctx context.Context, req *v1.GetReq) (*v1.GetResp, error) {
	server.mutex.Lock()
	v, ok := server.store[req.Key]
	server.mutex.Unlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Key not found %s",
			req.Key)
	}
	return &v1.GetResp{
		Val: v,
	}, nil
}
