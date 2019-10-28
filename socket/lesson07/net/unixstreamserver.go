package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func echoServer(c net.Conn) {
	defer func() {
		_ = c.Close()
	}()
	log.Printf("Client connected [%s]", c.RemoteAddr().Network())
	//io.Copy(c, c)
	var readLine [4096]byte
	for {
		n, err := c.Read(readLine[:])
		if n == 0 {
			log.Println("client quit")
			break
		}
		if err != nil {
			log.Fatal("Read error:", err)
		}
		log.Printf("Receive:%s\n", readLine[:n])
		sendLine := fmt.Sprintf("Hi, %s", readLine[:n])
		n, _ = c.Write([]byte(sendLine))
		if n != len(sendLine) {
			log.Fatal("write error:", err)
		}
	}
}

func main() {
	const SockAddr = "/tmp/echo.sock"
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("unix", SockAddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go echoServer(conn)
	}
}
