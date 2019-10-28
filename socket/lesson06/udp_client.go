package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"syscall"
)

func client(port int) {
	/* 创建字节流类型的IPV4 socket. */
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(sock)

	addr := syscall.SockaddrInet4{Port: port}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())

	scanner := bufio.NewScanner(os.Stdin)
	const MaxLine = 4096
	var recvLine [MaxLine]byte

	for scanner.Scan() {
		sendLine := scanner.Text()
		log.Println("now sending:", sendLine)
		err := syscall.Sendto(sock, []byte(sendLine), 0, &addr)
		if err != nil {
			log.Println(err)
		}
		// receive
		n, _, _ := syscall.Recvfrom(sock, recvLine[:], 0)
		recvLine[n] = '0'
		log.Println(string(recvLine[:]))
	}
}

func main() {
	client(12345)
}
