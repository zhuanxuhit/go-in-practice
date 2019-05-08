package main

import (
	"bufio"
	"flag"
	"github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	pb "github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/v3/tutorial"
	"github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/v3/util"
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
		//wr    = bufio.NewWriter(conn)
	)
	log.Infof("start tcp serve \"%s\" with \"%s\"", lAddr, rAddr)
	// 结构定义
	// 1. 读取头，MessageHeader 大小多少呢？
	messageHeader := &pb.MessageHeader{
		MessageType:   pb.Message_Message_XX1Request,
		MessageLength: 10,
	}
	out, err := proto.Marshal(messageHeader)
	if err != nil {
		log.Errorf("Failed to Marshal message:", err)
		return
	}
	headerSize := len(out)
	log.Infof("header size:%d", headerSize)
	request, err := util.ReadTCP(rr, headerSize)
	if err != nil {
		log.Errorf("Failed to ReadTCP message:", err)
		return
	}

	log.Infof("read request:%+v", request)

	defer conn.Close()
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
