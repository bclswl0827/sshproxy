package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/miekg/dns"
)

type DoHResolver struct {
	client *http.Client
	server string
	cache  sync.Map
}

func NewDoHResolver(server string) *DoHResolver {
	transport := &http.Transport{
		MaxIdleConns:       100,
		MaxConnsPerHost:    100,
		IdleConnTimeout:    90 * time.Second,
		DisableCompression: true,
		ForceAttemptHTTP2:  true,
	}
	return &DoHResolver{
		client: &http.Client{
			Timeout:   8 * time.Second,
			Transport: transport,
		},
		server: server,
	}
}

func (r *DoHResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	if ip := net.ParseIP(name); ip != nil {
		return ctx, ip, nil
	}

	if v, ok := r.cache.Load(name); ok {
		return ctx, v.(net.IP), nil
	}

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(name), dns.TypeA)

	raw, err := msg.Pack()
	if err != nil {
		return ctx, nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", r.server, bytes.NewReader(raw))
	if err != nil {
		return ctx, nil, err
	}
	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")

	resp, err := r.client.Do(req)
	if err != nil {
		return ctx, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx, nil, err
	}

	var dnsResp dns.Msg
	if err := dnsResp.Unpack(body); err != nil {
		return ctx, nil, err
	}
	if len(dnsResp.Answer) == 0 {
		return ctx, nil, errors.New("no DNS answer")
	}

	for _, ans := range dnsResp.Answer {
		if a, ok := ans.(*dns.A); ok {
			r.cache.Store(name, a.A)
			return ctx, a.A, nil
		}
	}
	return ctx, nil, errors.New("no A record found")
}
