package coreadblock

import (
	"bytes"
	"context"
	"testing"

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
}