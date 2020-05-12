package coreadblock

import (
	"fmt"
	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/ruijzhan/bloom"
	"sync"
)

func init()  {
	plugin.Register(PLUGIN_NAME, setup)
}


func setup(c *caddy.Controller) error  {
	a, err := adblockParse(c)
	if err != nil {
		log.Fatalf("%v", err)
	}

    go a.LoadRules()

	c.OnStartup(func() error {
		once.Do(func() {
			metrics.MustRegister(c, requestCount)
		})
		return nil
	})


	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		a.Next = next
		return a
	})

	return nil
}

func adblockParse(c *caddy.Controller) (*CoreAdBlock, error) {
	a := &CoreAdBlock{
		Exceptions: make(map[string]bool),
		Urls: []string{},
		BlockList: new(sync.Map),
		ResolveIP: "127.0.0.1",
	    Bloom: bloom.New(BLOOM_SIZE, HASH_SIZE)}

	for c.Next() {
		for c.NextBlock() {
    		p := c.Val()
    		switch p {
    		case "url":
    			a.Urls = append(a.Urls, c.RemainingArgs()...)
    		case "except":
    			for _, e := range c.RemainingArgs() {
    				a.Exceptions[e] = true
				}
    		case "ip":
    			a.ResolveIP = c.RemainingArgs()[0]
    		default:
    			return &CoreAdBlock{}, fmt.Errorf("unrecognized parameter %v", p)
    		}
		}
	}
	return a, nil
}