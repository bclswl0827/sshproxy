package main

import (
	"context"
	"log"
	"net"

	"github.com/armon/go-socks5"
	"golang.org/x/crypto/ssh"
)

func main() {
	var args Args
	args.Read()

	sshConn, err := ssh.Dial("tcp", args.Address, &ssh.ClientConfig{
		User: args.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(args.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatalln("failed to connect to server: ", err)
	}

	defer sshConn.Close()
	log.Println("successfully connected to server")

	serverSocks, err := socks5.New(&socks5.Config{
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return sshConn.Dial(network, addr)
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("starting SOCKS5 server at", args.Socks5)
	if err := serverSocks.ListenAndServe("tcp", args.Socks5); err != nil {
		log.Fatalln("failed to create SOCKS5 server", err)
	}
}
