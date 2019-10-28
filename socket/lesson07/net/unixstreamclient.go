package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func client(port int) {
	const SockAddr = "/tmp/echo.sock"
	conn, err := net.Dial("unix", SockAddr)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	const MaxLine = 4096
	var recvLine [MaxLine]byte

	for scanner.Scan() {
		sendLine := scanner.Bytes()
		log.Printf("now sending:%s\n", sendLine)
		n, err := conn.Write(sendLine)
		if n != len(sendLine) {
			log.Fatal("write error:", err)
		}
		n, _ = conn.Read(recvLine[:])
		log.Printf("read %s\n", recvLine[:n])
	}
}

func main() {
	client(12345)
}
