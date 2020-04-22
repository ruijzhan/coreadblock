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

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return CoreAdBlock{Next: next, Url: url}
	})

	return nil
	
}