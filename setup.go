package coreadblock

import (
	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
)

func init()  {
	plugin.Register(PLUGIN_NAME, setup)
}

func setup(c *caddy.Controller) error  {
	c.Next()

	if !c.NextArg(){
		return plugin.Error(PLUGIN_NAME, c.ArgErr())
	}

	url := c.Val()

	c.OnStartup(func() error {
		once.Do(func() {
			metrics.MustRegister(c, requestCount)
		})
		return nil
	})

	r, err := openURL(url)
	if err != nil {
		panic("Cannot read from url")
	}
	defer r.Close()
	m := newMap()
	m.parse(r)

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return CoreAdBlock{Next: next, Url: url, Data: m}
	})

	return nil
	
}