package main

import (
	"io"
	"log"
	"net"
)

func server(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()
	log.Printf("client connected:%s", conn.RemoteAddr())
	for {
		n, _ := io.Copy(conn, conn)
		if n == 0 {
			log.Printf("client:%s gone", conn.RemoteAddr())
			break
		}
		log.Printf("read %d bytes", n)
	}
}

//func main() {
//	l, err := net.Listen("tcp", ":12345")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	for {
//		conn, err := l.Accept()
//		if err != nil {
//			log.Fatalln(err)
//		}
//		go server(conn)
//	}
//}
