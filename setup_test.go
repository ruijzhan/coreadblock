package coreadblock

import (
	"testing"

	"github.com/caddyserver/caddy"
)

var corefile =
	` coreadblock {
        except www.qiudog.top www.163.com
        url https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts
        url https://adaway.org/hosts.txt
        ip 127.0.0.1
}
`

func TestSetup(t *testing.T) {
	c := caddy.NewTestController("dns", corefile)
	if err := setup(c); err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

//	c = caddy.NewTestController("dns", `coreadblock https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts`)
//	if err := setup(c); err == nil {
//		t.Fatalf("Expected errors, but got %v", err)
//	}
}

func TestAdblockParse(t *testing.T)  {
	c := caddy.NewTestController("dns", corefile)
	a, err := adblockParse(c)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if a.ResolveIP != "127.0.0.1" {
		t.Fatalf("Expected '127.0.0.1', but got %s", a.ResolveIP)
	}
	if len(a.Exceptions) != 2 {
		t.Fatalf("Expected 2, but got %d", len(a.Exceptions))
	}
	if len(a.Urls) != 2 {
		t.Fatalf("Expected 2, but got %d", len(a.Urls))
	}

}
