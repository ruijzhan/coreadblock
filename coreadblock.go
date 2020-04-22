package coreadblock

import (
	"context"
	"fmt"
	"github.com/coredns/coredns/plugin/metrics"
	"io"
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
}

func (ab CoreAdBlock) ServeDNS(cxt context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error)  {
	log.Debug("Received response")

	fmt.Println(r)

	requestCount.WithLabelValues(metrics.WithServer(cxt)).Inc()

	return ab.Next.ServeDNS(cxt, w, r)
}

func (_ CoreAdBlock) Name() string { return PLUGIN_NAME }
