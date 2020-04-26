package coreadblock

import (
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

func TestParse(t *testing.T)  {
	rc, err := openURL(url)
	if err != nil {
		t.Fatalf("Expected success, got failer")
	}
	defer rc.Close()
	m := newMap()
	m.parse(rc)
	t.Logf("Blocked list has %d entries", m.nBlocked)
}

