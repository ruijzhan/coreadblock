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
	"github.com/ruijzhan/bloom"
)

const (
	PLUGIN_NAME = "coreadblock"
	BLOOM_SIZE  = 50000 * 20
	HASH_SIZE   = 5
)

var (
	out io.Writer = os.Stdout
	log = clog.NewWithPlugin(PLUGIN_NAME)
)


type CoreAdBlock struct {
	Next 		plugin.Handler
	Urls		[]string
	ResolveIP   string
	Exceptions  map[string]bool
	BlockList   map[string]bool
	Bloom       *bloom.BloomFilter
	ready       bool
}

func (c *CoreAdBlock) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error)  {
	if !c.ready {
		log.Info("Plugin not Ready")
		return plugin.NextOrFailure(c.Name(), c.Next, ctx, w, r)
	}

	state := request.Request{W:w, Req: r}
	qname := state.Name()
	log.Infof("%d entries in blacklist", len(c.BlockList))
	log.Infof("Received qname %d", qname)

	var answers []dns.RR

	if state.QType() == dns.TypeA {
		if c.Exceptions[qname] {
			// do nothing
		} else if c.Bloom.Test([]byte(qname)) {
			log.Infof("Filtered %s", qname)
			ips := []net.IP{net.ParseIP(c.ResolveIP)}
			answers = a(qname, 3600, ips)
		}
	}

	if len(answers) == 0 {
		return plugin.NextOrFailure(c.Name(), c.Next, ctx, w, r)
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = answers
	w.WriteMsg(m)

	requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

	return dns.RcodeSuccess, nil
}

func (c *CoreAdBlock) Name() string { return PLUGIN_NAME }

func (c *CoreAdBlock) LoadRules() error {
	for _, url := range c.Urls {
		if err := c.parseHostsURL(url); err != nil{
			log.Warningf("Failed to parse url %v because %v", url, err )
			return err
		}
	}
	c.ready = true
	return nil
}

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
