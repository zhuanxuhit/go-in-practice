package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// 读取数据
	scanner := bufio.NewScanner(os.Stdin)
	reader := bufio.NewReader(conn)

	for scanner.Scan() {
		input := scanner.Text()
		fmt.Printf("read from stdin:%s\n", input)

		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			panic(err)
		}
		output, err := reader.ReadString('\n') // 此处末尾的\n不会去除
		if err != nil {
			panic(err)
		}
		fmt.Printf("read from remote:%s", output)
	}
	if scanner.Err() != nil {
		fmt.Printf("read scanner error:%s\n", scanner.Err())
	}
}
