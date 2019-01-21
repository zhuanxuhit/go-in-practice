package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

type SSHtunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

	Config *ssh.ClientConfig
}

func (tunnel *SSHtunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHtunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		fmt.Printf("Server dial error: %s\n", err)
		return
	}

	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func expandUser(s string) string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	home := os.Getenv(env)
	if home == "" {
		return s
	}

	if len(s) >= 2 && s[0] == '~' && os.IsPathSeparator(s[1]) {
		if runtime.GOOS == "windows" {
			s = filepath.ToSlash(filepath.Join(home, s[2:]))
		} else {
			s = filepath.Join(home, s[2:])
		}
	}
	return os.Expand(s, func(env string) string {
		if env == "HOME" {
			return home
		}
		return os.Getenv(env)
	})
}

func SSHPublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("read %s error. ", file)
		return nil
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		fmt.Printf("ParsePrivateKey %s error. ", file)
		return nil
	}
	return ssh.PublicKeys(key)
}
func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func main() {
	localEndpoint := Endpoint{
		Host: "localhost",
		Port: 9000,
	}

	serverEndpoint := Endpoint{
		Host: "10.8.122.214",
		Port: 22,
	}

	remoteEndpoint := &Endpoint{
		Host: "localhost",
		Port: 8080,
	}

	sshConfig := ssh.ClientConfig{
		User: "wangchao.zhuanxu",
		Auth: []ssh.AuthMethod{
			//SSHAgent(),
			SSHPublicKeyFile(expandUser("~/.ssh/id_rsa")),
		},
		//需要验证服务端，不做验证返回nil就可以，点击HostKeyCallback看源码就知道了
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			//fmt.Println(hostname, remote, key)
			return nil
		},
	}

	tunnel := &SSHtunnel{
		Config: &sshConfig,
		Local:  &localEndpoint,
		Server: &serverEndpoint,
		Remote: remoteEndpoint,
	}

	tunnel.Start()
}
