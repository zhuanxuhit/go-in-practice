package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func SSHPassword() ssh.AuthMethod {
	return ssh.Password("123456")
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

type SSHClient struct {
	Config *ssh.ClientConfig
	Host   string
	Port   int
	client *ssh.Client
}

type SSHCommand struct {
	Path   string
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (client *SSHClient) RunCommand(cmd *SSHCommand) error {
	var (
		session *ssh.Session
		err     error
	)

	if session, err = client.newSession(); err != nil {
		return err
	}
	defer session.Close()

	if err = client.prepareCommand(session, cmd); err != nil {
		return err
	}

	err = session.Run(cmd.Path)
	return err
}
func (client *SSHClient) close() {
	if client.client != nil {
		client.client.Close()
	}
}

func (client *SSHClient) prepareCommand(session *ssh.Session, cmd *SSHCommand) error {
	for _, env := range cmd.Env {
		variable := strings.Split(env, "=")
		if len(variable) != 2 {
			continue
		}

		if err := session.Setenv(variable[0], variable[1]); err != nil {
			return err
		}
	}

	if cmd.Stdin != nil {
		stdin, err := session.StdinPipe()
		if err != nil {
			return fmt.Errorf("unable to setup stdin for session: %v", err)
		}
		go io.Copy(stdin, cmd.Stdin)
	}

	if cmd.Stdout != nil {
		stdout, err := session.StdoutPipe()
		if err != nil {
			return fmt.Errorf("unable to setup stdout for session: %v", err)
		}
		go io.Copy(cmd.Stdout, stdout)
	}

	if cmd.Stderr != nil {
		stderr, err := session.StderrPipe()
		if err != nil {
			return fmt.Errorf("unable to setup stderr for session: %v", err)
		}
		go io.Copy(cmd.Stderr, stderr)
	}

	return nil
}

func (client *SSHClient) newSession() (*ssh.Session, error) {
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port), client.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %s", err)
	}
	client.client = connection

	session, err := connection.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %s", err)
	}

	modes := ssh.TerminalModes{
		// ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	return session, nil
}

func main() {
	// 1. 生成 ssh 配置

	//fmt.Println(expandUser("~/.ssh/id_rsa"))
	//return
	sshConfig := ssh.ClientConfig{
		User: "wangchao.zhuanxu",
		Auth: []ssh.AuthMethod{
			//SSHPassword(),
			SSHPublicKeyFile(expandUser("~/.ssh/id_rsa")),
		},
		//需要验证服务端，不做验证返回nil就可以，点击HostKeyCallback看源码就知道了
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			fmt.Println(hostname, remote, key)
			return nil
		},
	}
	// 2. 生成 ssh 客户端
	client := &SSHClient{
		Config: &sshConfig,
		Host:   "10.8.122.214",
		Port:   22,
	}
	defer client.close()
	// 3. 定义cmd
	cmd := SSHCommand{
		Path:   "ls -l $LC_DIR",
		Env:    []string{"LC_DIR=/"},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	fmt.Printf("Running command: %s\n", cmd.Path)
	if err := client.RunCommand(&cmd); err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
		os.Exit(1)
	}
	time.Sleep(1 * time.Second)
}
