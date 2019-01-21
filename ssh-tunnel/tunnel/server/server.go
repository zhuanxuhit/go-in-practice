package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	for {
		l, err := conn.Accept()
		if err != nil {
			panic(err)
		}
		go func(conn net.Conn) {
			defer conn.Close()
			// 读取数据
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				input := scanner.Text()
				fmt.Printf("read from remote:%s\n", input)

				_, err := conn.Write([]byte(input+"\n"))
				if err != nil {
					panic(err)
				}
			}
			if scanner.Err() != nil {
				fmt.Printf("read scanner error:%s\n", scanner.Err())
			}
		}(l)
	}
}
