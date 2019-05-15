package v1

import (
	"context"
	"github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/cache/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
