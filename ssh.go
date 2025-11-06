package main

import (
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHDialer struct {
	addr     string
	username string
	password string

	mu   sync.Mutex
	conn *ssh.Client
}

func NewSSHDialer(addr, username, password string) *SSHDialer {
	return &SSHDialer{
		addr:     addr,
		username: username,
		password: password,
	}
}

func (r *SSHDialer) ensureConn() (*ssh.Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.conn != nil {
		// keepalive test
		_, _, err := r.conn.SendRequest("keepalive@openssh.com", true, nil)
		if err == nil {
			return r.conn, nil
		}
		log.Println("[SSH] Connection lost:", err)
		r.conn.Close()
		r.conn = nil
	}

	for {
		log.Println("[SSH] Connecting to", r.addr)
		c, err := ssh.Dial("tcp", r.addr, &ssh.ClientConfig{
			User:            r.username,
			Auth:            []ssh.AuthMethod{ssh.Password(r.password)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Config: ssh.Config{
				MACs: []string{"hmac-sha2-256", "hmac-sha2-512", "hmac-sha1"},
			},
			Timeout: 10 * time.Second,
		})
		if err == nil {
			r.conn = c
			log.Println("[SSH] Connected successfully")
			return c, nil
		}
		log.Println("[SSH] Connect failed:", err)
		time.Sleep(5 * time.Second)
	}
}

func (r *SSHDialer) Dial(network, addr string) (net.Conn, error) {
	c, err := r.ensureConn()
	if err != nil {
		return nil, err
	}
	conn, err := c.Dial(network, addr)
	if err != nil {
		log.Println("[SSH] Dial failed, retrying:", err)
		r.mu.Lock()
		r.conn.Close()
		r.conn = nil
		r.mu.Unlock()
		return r.Dial(network, addr)
	}
	return conn, nil
}
