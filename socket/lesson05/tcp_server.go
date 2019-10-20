package main

import (
	"flag"
	"syscall"
	"log"
	"net"
	"os"
	"time"
)

var port int

func init() {
	flag.IntVar(&port, "port", 12345, "port")
}

func readn(fd int, buf []byte, size int) (int, error) {
	length := size
	bufPoint := buf
	for length > 0 {
		result, err := syscall.Read(fd, bufPoint)
		if err != nil {
			if err == syscall.EINTR {
				continue
			} else {
				return -1, err
			}
		}
		if result == 0 {
			break
		}
		length -= result
		bufPoint = bufPoint[result:]
	}
	return size - length, nil
}

func readData(sockfd int) {
	var buf [1024]byte
	times := 0
	for {
		log.Printf("block in read\n")
		n, err := readn(sockfd, buf[0:], 1024)
		if n == 0 {
			return
		}
		if err != nil {
			log.Printf("read err:%v\n", err)
		} else {
			times += 1
			log.Printf("1K read for %d \n", times)
		}
		time.Sleep(1000 * time.Microsecond)
	}
}

func server(port int) {
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
	/* listen的backlog为1024 */
	err = syscall.Listen(sock, 1024)
	if err != nil {
		log.Fatal(err)
	}
	for {
		connfd, cliaddr, err := syscall.Accept(sock)
		if addr, ok := cliaddr.(*syscall.SockaddrInet4); ok {
			log.Printf("client: %s:%d\n", net.IPAddr{IP: addr.Addr[0:]}, addr.Port)
		}
		if err != nil {
			log.Println("accept: ", err)
			continue
		}
		readData(connfd)
		syscall.Close(connfd)
	}
}

//https://github.com/froghui/yolanda/chap-4/make_socket.c
func main() {
	flag.Parse()
	log.Printf("listen to port:%d\n", port)
	server(port)
	os.Exit(0)
}
