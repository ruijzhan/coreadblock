package coreadblock

import (
	"github.com/caddyserver/caddy"
	"testing"
)

var url = "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts"

func TestOpenURL(t *testing.T) {
	rc, err := openURL(url)
	if err != nil {
		t.Fatalf("Expected success, got failer")
	}
	defer rc.Close()
}

func TestParseHosts(t *testing.T)  {
	c := caddy.NewTestController("dns", corefile)
	adblk, err := adblockParse(c)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	for _, url := range adblk.Urls{
		if err = adblk.parseHostsURL(url); err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
	}
	t.Logf("Parsed %d entries", len(adblk.BlockList))
	if !adblk.BlockList["cdn.3lift.com"] {
		t.Fatalf("Expected cdn.3lift.com in blocked list, but no")
	}
	if !adblk.Exceptions["www.qiudog.top"] {
		t.Fatalf("Expected www.qiudog.top in exception list, but no")
	}
}

