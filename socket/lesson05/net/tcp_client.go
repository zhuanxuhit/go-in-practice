package main

import (
	"log"
	"net"
)


func writeAll(conn net.Conn, buf []byte) bool {
	for len(buf) > 0 {
		// /usr/local/opt/go/libexec/src/net/net.go
		// /usr/local/opt/go/libexec/src/net/tcpsock.go
		n, err := conn.Write(buf)
		if err != nil {
			log.Println("send error:", err)
			return false
		}
		log.Println("send into buffer ", n)
		buf = buf[n:]
	}
	return true
}

func sendData(conn net.Conn) {
	const MessageSize = 1024
	var query [MessageSize + 1]byte
	for i := 0; i < MessageSize; i++ {
		query[i] = 'a'
	}
	query[MessageSize] = '0'
	writeAll(conn, query[0:])
}

func main() {
	// /usr/local/opt/go/libexec/src/net/tcpsock_posix.go
	conn, err := net.Dial("tcp", "0.0.0.0:12345")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	sendData(conn)
}
