package coreadblock

import (
	"bytes"
	"context"
	"testing"
	"net"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"

	"github.com/miekg/dns"
)

func TestCoreAdBlock(t *testing.T) {
	adblk := CoreAdBlock{Next: test.ErrorHandler(), Url: `https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts`}

	b := &bytes.Buffer{}
	out = b

	ctx := context.TODO()
	r := new(dns.Msg)
	r.SetQuestion("example.org", dns.TypeA)

	rec := dnstest.NewRecorder(&test.ResponseWriter{})
	adblk.ServeDNS(ctx, rec, r)
	for _, rr := range rec.Msg.Answer {
		r:= rr.(*dns.A)
		if !r.A.Equal(net.ParseIP("127.0.0.1")) {
			t.Fatalf("Expected 127.0.0.1, but got %v", r.A)
		}
	}
}