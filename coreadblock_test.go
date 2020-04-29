package coreadblock


import (
	"bytes"
	"context"
	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	"math/rand"
	"net"
	"strings"
	"testing"
)

func TestCoreAdBlock(t *testing.T) {
	c := caddy.NewTestController("dns", corefile)
	a, err := adblockParse(c)
    a.parseHosts(strings.NewReader(hostsSample55k))
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	b := &bytes.Buffer{}
	out = b

	ctx := context.TODO()
	r := new(dns.Msg)
	r.SetQuestion("cdn.3lift.com", dns.TypeA)

	rec := dnstest.NewRecorder(&test.ResponseWriter{})
	a.ServeDNS(ctx, rec, r)
	if rec.Rcode == dns.RcodeSuccess {
		for _, rr := range rec.Msg.Answer {
			r:= rr.(*dns.A)
			if !r.A.Equal(net.ParseIP(a.ResolveIP)) {
				t.Fatalf("Expected 127.0.0.1, but got %v", r.A)
			}
		}
	}
}

func BenchmarkResolve1k(b *testing.B){
	benchmarkResolv(b, hostsSample1k)
}

func BenchmarkResolve55k(b *testing.B){
	benchmarkResolv(b, hostsSample55k)
}

func benchmarkResolv(b *testing.B, hosts string){
	c := caddy.NewTestController("dns", corefile)
	a, err := adblockParse(c)
	a.parseHosts(strings.NewReader(hosts))
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
