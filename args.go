package main

import "flag"

func (a *Args) Read() {
	flag.StringVar(&a.Address, "address", "127.0.0.1:22", "SSH server address (host:port)")
	flag.StringVar(&a.Username, "username", "root", "SSH username")
	flag.StringVar(&a.Password, "password", "passw0rd", "SSH password")
	flag.StringVar(&a.Socks5, "socks5", "0.0.0.0:10808", "SOCKS5 server listen address (host:port)")
	flag.StringVar(&a.DoH, "doh", "https://120.53.53.53/dns-query", "DoH server address")
	flag.Parse()
}
