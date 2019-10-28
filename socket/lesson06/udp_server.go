package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func server(port int) {
	/* 创建字节流类型的IPV4 socket. */
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
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
	const MaxLine = 4096
	var message [MaxLine]byte
	var sendLine bytes.Buffer

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		for _ = range sigs {
			;
		}
	}()

	for {

		n, cliaddr, _ := syscall.Recvfrom(sock, message[:], 0)
		if addr, ok := cliaddr.(*syscall.SockaddrInet4); ok {
			log.Printf("client: %s:%d\n", net.IPAddr{IP: addr.Addr[0:]}, addr.Port)
		}
		//message[n] = '0'
		log.Printf("received %d bytes: %s\n", n, string(message[:n]))

		//buf := bytes.NewBuffer(sendLine[:])
		sendLine.Reset()
		_, _ = fmt.Fprintf(&sendLine, "Hi, %s", message[:n])
		err := syscall.Sendto(sock, sendLine.Bytes(), 0, cliaddr)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	server(12345)
}
