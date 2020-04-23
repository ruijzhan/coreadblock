package coreadblock

import (
	"testing"

	"github.com/caddyserver/caddy"
)

func TestSetup(t *testing.T) {
	c := caddy.NewTestController("dns", `coreadblock https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts`)
	if err := setup(c); err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

//	c = caddy.NewTestController("dns", `coreadblock https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts`)
//	if err := setup(c); err == nil {
//		t.Fatalf("Expected errors, but got %v", err)
//	}
}
