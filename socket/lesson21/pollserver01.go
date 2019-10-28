package main

import (
	"golang.org/x/sys/unix"
	"log"
	"net"
	"syscall"
)

const InitSize = 1024

func main() {
	/* 创建字节流类型的IPV4 socket. */
	sock, err := unix.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatalln(err)
	}

	addr := unix.SockaddrInet4{Port: 12345}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
	/* 绑定到port和ip. */
	err = unix.Bind(sock, &addr)
	if err != nil {
		log.Fatalln(err)
	}
	/* listen的backlog为1024 */
	err = unix.Listen(sock, 1024)
	if err != nil {
		log.Fatal(err)
	}
	eventSet := make([]unix.PollFd, 0)
	eventSet = append(eventSet, unix.PollFd{
		Fd:     int32(sock),
		Events: unix.EPOLLRDNORM,
	})

	defer func() {
		for _, event := range eventSet {
			if event.Fd != -1 {
				_ = unix.Close(int(event.Fd))
			}
		}
	}()

	var buf [1024]byte
	for {
		readyNumber, err := unix.Poll(eventSet, -1)
		if err != nil {
			log.Printf("poll err: %v", err)
			return
		}
		if (eventSet[0].Revents & unix.EPOLLRDNORM) > 0 {
			connfd, cliaddr, err := syscall.Accept(sock)
			if err != nil {
				log.Printf("Accept err: %v", err)
				return
			}
			if addr, ok := cliaddr.(*syscall.SockaddrInet4); ok {
				log.Printf("client: %s:%d\n", net.IPAddr{IP: addr.Addr[0:]}, addr.Port)
			}
			notFound := true
			for i := 1; i < len(eventSet); i++ {
				if eventSet[i].Fd != -1 {
					eventSet[i].Fd = int32(connfd)
					eventSet[i].Events = unix.EPOLLRDNORM
					notFound = false
					break
				}
			}
			if notFound {
				eventSet = append(eventSet, unix.PollFd{
					Fd:     int32(connfd),
					Events: unix.EPOLLRDNORM,
				})
			}
			readyNumber -= 1
			if readyNumber <= 0 {
				continue
			}
		}
		for i := range eventSet {
			event := &eventSet[i]
			if event.Fd == -1 {
				continue
			}

			if (event.Revents & (unix.EPOLLRDNORM | unix.POLLERR)) > 0 {
				n, _ := unix.Read(int(event.Fd), buf[:])
				log.Printf("read %d byte, str:%s\n", n, buf[:n])
				if n == 0 {
					unix.Close(int(event.Fd))
					event.Fd = -1
				} else if n > 0 {
					_, _ = unix.Write(int(event.Fd), buf[:n])
				}
				readyNumber -= 1
				if readyNumber <= 0 {
					continue
				}
			}
		}
	}
}
