package main

import (
	"flag"
	log "github.com/golang/glog"
	"github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/v2/method"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
)



func main() {
	flag.Parse()
	arith := new(method.Arith)
	//var DefaultServer = rpc.NewServer()
	// func Register(rcvr interface{}) error { return DefaultServer.Register(rcvr) }
	rpc.Register(arith)
	// DefaultServer.HandleHTTP(DefaultRPCPath, DefaultDebugPath)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("my-rpc-srv get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
