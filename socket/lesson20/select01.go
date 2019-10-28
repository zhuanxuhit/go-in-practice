package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"syscall"
)

func main() {
	/* 创建字节流类型的IPV4 socket. */
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(sock)

	addr := syscall.SockaddrInet4{Port: 12345}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
	/* 连接到server. */
	err = syscall.Connect(sock, &addr)
	if err != nil {
		log.Fatalln(err)
	}

	readmask, allread := syscall.FdSet{}, syscall.FdSet{}
	//timeout := &syscall.Timeval{}
	FD_ZERO(&allread)
	FD_SET(&allread, sock)
	FD_SET(&allread, int(os.Stdin.Fd()))
	var buf [1024]byte

	for {
		readmask = allread

		_, err = syscall.Select(sock+1, &readmask, nil, nil, nil)
		if err != nil {
			log.Fatalln(err)
		}
		if FD_ISSET(&readmask, sock) {
			// 读数据
			n, err := syscall.Read(sock, buf[:])
			if err != nil {
				log.Fatalf("read error:%v\n", err)
			}

			if n == 0 {
				log.Println("server has gone.")
				break
			}
			log.Printf("%s\n", buf[:n])
		}
		if FD_ISSET(&readmask, int(os.Stdin.Fd())) {
			consoleReader := bufio.NewReader(os.Stdin)
			b, err := consoleReader.ReadBytes('\n')
			//b, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("read error:%v\n", err)
			}
			//b[len(b)-1] = '0'
			log.Printf("now sending:%s", string(b[:len(b)-1]))
			n, err := syscall.Write(sock, b[:len(b)-1])
			if err != nil {
				log.Fatalf("write error:%v\n", err)
			}
			log.Printf("send byte:%d", n)
		}
	}
}

func FD_SET(p *syscall.FdSet, i int) {
	p.Bits[i/32] |= 1 << (uint(i) % 32)
}

func FD_ISSET(p *syscall.FdSet, i int) bool {
	return (p.Bits[i/32] & (1 << (uint(i) % 32))) != 0
}

func FD_ZERO(p *syscall.FdSet) {
	for i := range p.Bits {
		p.Bits[i] = 0
	}
}
