package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/golang/glog"
	"net"
	"net/textproto"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func serveTCP(conn *net.TCPConn) {
	var (
		lAddr = conn.LocalAddr().String()
		rAddr = conn.RemoteAddr().String()
		rr    = bufio.NewReader(conn)
		line  string
		err   error
		wr    = bufio.NewWriter(conn)
	)
	defer conn.Close()

	log.Infof("start tcp serve \"%s\" with \"%s\"", lAddr, rAddr)
	// 1.读取文件头
	//buf, _ := ioutil.ReadAll(rr)
	//log.Infof("read data:%q", buf)
	tp := textproto.NewReader(rr)
	// First line: GET /index.html HTTP/1.0
	if line, err = tp.ReadLine(); err != nil {
		log.Errorf("读取请求行失败:%s", err)
		return
	}
	// 解析出 method, requestURI, proto string
	var method, requestURI, proto string
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	s2 += s1 + 1
	method, requestURI, proto = line[:s1], line[s1+1 : s2], line[s2+1:]
	log.Infof("method:%s, requestURI:%s, proto:%s", method, requestURI, proto)
	// 2. 读取首部
	s := rr.Buffered()
	peek, _ := rr.Peek(s)
	var header []string
	for len(peek) > 0 {
		i := bytes.IndexByte(peek, '\n')
		item := string(peek[:i-1])
		if len(item) == 0 {
			peek = peek[i+1:] // 推进内容
			break             // \r\n\r\n 连续两个，表示头部读完
		}
		header = append(header, item) // 去除 \r
		peek = peek[i+1:]
	}
	log.Infof("\nheader is:%v", strings.Join(header, "\n"))
	log.Infof("body is:%q", peek)

	responseBegin := fmt.Sprintf("%s %d %s", "HTTP/1.1", 200, "ok")
	responseHeader := []string{
		"content-type: application/json",
	}
	msg := map[string]string{
		"begin": line,
		//"header": strings.Join(header, "\n"),
		"body": string(peek),
	}
	data, _ := json.Marshal(msg)
	responseHeader = append(responseHeader, fmt.Sprintf("content-length: %d", len(data)))
	response := responseBegin + "\r\n" + strings.Join(responseHeader, "\r\n") + "\r\n\r\n" + string(data)

	wr.Write([]byte(response))
	wr.Flush()
}
func acceptTCP(lis *net.TCPListener) {
	for {
		conn, err := lis.AcceptTCP()
		if err != nil {
			log.Infof("listener.Accept(\"%s\") error(%v)", lis.Addr().String(), err)
			return
		}
		go serveTCP(conn)
	}
	//http.Serve()
}

func main() {
	flag.Parse()

	addr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		log.Fatalln("ResolveTCPAddr error:", err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln("network error:", err)
	}

	go acceptTCP(l)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("my-http-srv get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
