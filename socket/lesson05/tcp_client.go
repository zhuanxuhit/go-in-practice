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

func sendData(fd int) {
	const MessageSize = 1024
	var query [MessageSize + 1]byte
	for i := 0; i < MessageSize; i++ {
		query[i] = 'a'
	}
	query[MessageSize] = '0'

	remaining := len(query)
	cp := query[0:]
	for remaining > 0 {
		nWritten, err := syscall.Write(fd, cp)
		if err != nil {
			log.Println("send error:", err)
			return
		}
		log.Println("send into buffer ", nWritten)
		remaining -= nWritten
		cp = cp[nWritten:]
	}
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
