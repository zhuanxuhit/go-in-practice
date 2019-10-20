package main

import (
	"flag"
	"log"
	"os"
	"syscall"
	"net"
)

var port int

func init() {
	flag.IntVar(&port, "port", 12345, "port")
}

func writeAll(fd int, buf []byte) bool {
	for len(buf) > 0 {
		n, err := syscall.Write(fd, buf)
		if err != nil {
			log.Println("send error:", err)
			return false
		}
		log.Println("send into buffer ", n)
		buf = buf[n:]
	}
	return true
}

func sendData(fd int) {
	const MessageSize = 1024
	var query [MessageSize + 1]byte
	for i := 0; i < MessageSize; i++ {
		query[i] = 'a'
	}
	query[MessageSize] = '0'
	writeAll(fd, query[0:])
}

func client(port int) {
	/* 创建字节流类型的IPV4 socket. */
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(sock)

	addr := syscall.SockaddrInet4{Port: port}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
	/* 连接到server. */
	err = syscall.Connect(sock, &addr)
	if err != nil {
		log.Fatal(err)
	}
	sendData(sock)
}

func main() {
	flag.Parse()
	log.Printf("connect to port:%d\n", port)
	client(port)
	os.Exit(0)
}
