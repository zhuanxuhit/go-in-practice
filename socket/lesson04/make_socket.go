package main

import (
	"flag"
	"os"
	"syscall"
	"log"
	"net"
)

var port int

func init() {
	flag.IntVar(&port, "port", 12345, "port")
}

func make_socket(port int) {
	/* 创建字节流类型的IPV4 socket. */
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(sock)

	addr := syscall.SockaddrInet4{Port: port}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
	/* 绑定到port和ip. */
	err = syscall.Bind(sock, &addr)
	if err != nil {
		log.Fatal(err)
	}
	err = syscall.Listen(sock, 10)
	if err != nil {
		log.Fatal(err)
	}
	for {
		;
	}

}

//https://github.com/froghui/yolanda/chap-4/make_socket.c
func main() {
	flag.Parse()
	log.Printf("listen to port:%d\n", port)
	make_socket(port)
	os.Exit(0)
}
