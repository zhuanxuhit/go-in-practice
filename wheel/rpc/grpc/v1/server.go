package main

import (
	"bufio"
	"flag"
	log "github.com/golang/glog"
	pb "github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/v1/tutorial"
	"github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/v1/util"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func serveTCP(conn *net.TCPConn) {
	var (
		lAddr = conn.LocalAddr().String()
		rAddr = conn.RemoteAddr().String()
		rr    = bufio.NewReader(conn)
		wr    = bufio.NewWriter(conn)
	)
	log.Infof("start tcp serve \"%s\" with \"%s\"", lAddr, rAddr)
	// 结构定义
	// 1. 读取头，一字节，表示长度
	message, err := util.ReadTCP(rr)
	if err != nil {
		log.Errorf("readTCP err:%v", err)
		return
	}
	log.Info("read message:", message)
	// 3. 根据不同的消息类型，返回不同的消息响应
	if message.Request.Request1 != nil {
		message.Response = &pb.Response{
			Response1: &pb.XX1Response{
				Name: "XX1Response",
			},
		}
	} else if message.Request.Request2 != nil {
		message.Response = &pb.Response{
			Response2: &pb.XX2Response{
				Name: "XX2Response",
			},
		}
	}
	message.Request = nil
	// 开始写返回值
	err = util.WriteTCP(wr, message)
	if err != nil {
		log.Errorf("writeTCP err:%v", err)
		return
	}
}
func acceptTCP(lis *net.TCPListener) {
	for {
		conn, err := lis.AcceptTCP()
		if err != nil {
			log.Errorf("listener.Accept(\"%s\") error(%v)", lis.Addr().String(), err)
			return
		}
		go serveTCP(conn)
	}
}

func main() {
	flag.Parse()
	addr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		log.Fatalln("ResolveTCPAddr error:", err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln("network error:", err)
	}
	go acceptTCP(l)

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
