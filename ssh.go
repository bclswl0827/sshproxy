package main

import (
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

type SSHDialer struct {
	client *ssh.Client
	pool   sync.Pool
}

func NewSSHDialer(client *ssh.Client) *SSHDialer {
	return &SSHDialer{
		client: client,
		pool: sync.Pool{
			New: func() any {
				return client
			},
		},
	}
}

func (d *SSHDialer) Dial(network, addr string) (net.Conn, error) {
	c := d.pool.Get().(*ssh.Client)
	conn, err := c.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
