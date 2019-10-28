package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 12345})
	if err != nil {
		log.Fatal(err)
	}
	const MaxLine = 4096
	var message [MaxLine]byte

	for {
		n, remoteAddr, err := listener.ReadFromUDP(message[:])
		if err != nil {
			log.Printf("error during read: %s", err)
		}
		log.Printf("received %d bytes: %s\n", n, string(message[:n]))

		_, err = listener.WriteToUDP([]byte(fmt.Sprintf("Hi, %s", message[:n])), remoteAddr)
		if err != nil {
			fmt.Printf(err.Error())
		}
	}
}
