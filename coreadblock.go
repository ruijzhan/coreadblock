package coreadblock

import (
	"context"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/coredns/coredns/request"
	"io"
	"net"
	"os"

	"github.com/coredns/coredns/plugin"

	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
)

const PLUGIN_NAME = "coreadblock"

var (
	out io.Writer = os.Stdout
	log = clog.NewWithPlugin(PLUGIN_NAME)
)


type CoreAdBlock struct {
	Next plugin.Handler
	Url	 string
	Data *Map
}

func (ab CoreAdBlock) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error)  {
	state := request.Request{W:w, Req: r}
	qname := state.Name()

	var answers []dns.RR

	switch state.QType() {
	case dns.TypeA:
		ips := []net.IP{net.ParseIP("127.0.0.1")}
		answers = a(qname, 3600, ips)
	}

	if len(answers) == 0 {
		return plugin.NextOrFailure(ab.Name(), ab.Next, ctx, w, r)
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = answers

	w.WriteMsg(m)

	requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

	return dns.RcodeSuccess, nil
}

func (_ CoreAdBlock) Name() string { return PLUGIN_NAME }

func a(zone string, ttl uint32, ips []net.IP) []dns.RR {
	answers := make([]dns.RR, len(ips))
	for i, ip := range ips {
		r := new(dns.A)
		r.Hdr = dns.RR_Header{Name: zone, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: ttl}
		r.A = ip
		answers[i] = r
	}
	return answers
}
