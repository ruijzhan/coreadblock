package coreadblock

import (
	"context"
	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	"io"
	"math/rand"
	"net"
	"testing"
)

func TestCoreAdBlock(t *testing.T) {
	rc := hostsSample1k()
	defer rc.Close()
	c := caddy.NewTestController("dns", corefile)
	a, err := adblockParse(c)
	a.Ready = true
    a.parseHosts(rc)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	keys := make([]string, 0, len(a.BlockList))
	for k, _ := range a.BlockList {
		keys = append(keys, k)
	}

	ctx := context.TODO()
	for _, k := range keys {
		r := new(dns.Msg)
		r.SetQuestion(k, dns.TypeA)

		rec := dnstest.NewRecorder(&test.ResponseWriter{})
		a.ServeDNS(ctx, rec, r)
		if rec.Rcode == dns.RcodeSuccess && rec.Msg != nil{
			for _, rr := range rec.Msg.Answer {
				r:= rr.(*dns.A)
				if !r.A.Equal(net.ParseIP(a.ResolveIP)) {
					t.Fatalf("Expected 127.0.0.1, but got %v", r.A)
				}
			}
		}
	}
}

func BenchmarkResolve1k(b *testing.B){
	benchmarkResolv(b, hostsSample1k())
}

func BenchmarkResolve55k(b *testing.B){
	benchmarkResolv(b, hostsSample55k())
}

func benchmarkResolv(b *testing.B, hosts io.ReadCloser){
	defer hosts.Close()
	c := caddy.NewTestController("dns", corefile)
	a, err := adblockParse(c)
	a.parseHosts(hosts)
	a.Ready = true
	if err != nil {
		b.Fatalf("Expected no error, but got %v", err)
	}

	keys := make([]string, len(a.BlockList))
	for k, _ := range a.BlockList {
		keys = append(keys, k)
	}

	ctx := context.TODO()
	r := new(dns.Msg)

	rec := dnstest.NewRecorder(&test.ResponseWriter{})
	for i :=0; i<b.N; i++ {
		r.SetQuestion(keys[rand.Intn(len(keys))], dns.TypeA)
		a.ServeDNS(ctx, rec, r)
	}
}
