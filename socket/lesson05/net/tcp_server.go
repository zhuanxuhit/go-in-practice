package main

import (
	"log"
	"net"
	"time"
)

func readn(conn net.Conn, buf []byte, size int) (int, error) {
	length := size
	bufPoint := buf
	for length > 0 {
		// /usr/local/opt/go/libexec/src/internal/poll/fd_unix.go
		// 里面处理了 EAGAIN 和 EINTR error
		result, err := conn.Read(bufPoint)
		if result == 0 {
			break
		}
		if err != nil {
			return -1, err
		}
		length -= result
		bufPoint = bufPoint[result:]
	}
	return size - length, nil
}

func readData(conn net.Conn) {
	var buf [1024]byte
	times := 0
	for {
		log.Printf("block in read\n")
		// /usr/local/opt/go/libexec/src/internal/poll/fd_unix.go
		n, err := readn(conn, buf[:], len(buf))
		if n == 0 {
			return
		}
		if err != nil {
			log.Printf("read err:%v\n", err)
		} else {
			times += 1
			log.Printf("1K read for %d, read bytes:%d\n", times, n)
		}
		time.Sleep(1000 * time.Microsecond)
	}
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:12345")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error: ", err)
			continue
		}
		readData(conn)
		conn.Close()
	}
}
