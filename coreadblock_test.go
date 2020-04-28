package coreadblock


import (
	"bytes"
	"context"
	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	"net"
	"testing"
)

func TestCoreAdBlock(t *testing.T) {
	c := caddy.NewTestController("dns", corefile)
	adblk, err := adblockParse(c)
	for _, url := range adblk.Urls {
		if err := adblk.parseHostsURL(url); err != nil{
			log.Warningf("Failed to parse url %v because %v", url, err )
		}
	}
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	b := &bytes.Buffer{}
	out = b

	ctx := context.TODO()
	r := new(dns.Msg)
	r.SetQuestion("cdn.3lift.com", dns.TypeA)

	rec := dnstest.NewRecorder(&test.ResponseWriter{})
	adblk.ServeDNS(ctx, rec, r)
	if rec.Rcode == dns.RcodeSuccess {
		for _, rr := range rec.Msg.Answer {
			r:= rr.(*dns.A)
			if !r.A.Equal(net.ParseIP(adblk.ResolveIP)) {
				t.Fatalf("Expected 127.0.0.1, but got %v", r.A)
			}
		}
	}
}