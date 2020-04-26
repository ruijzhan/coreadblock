package coreadblock

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"
)

func parseIP(addr string) net.IP  {
	if i := strings.Index(addr, "%"); i >= 0 {
		// discard ipv6 zone
		addr = addr[:i]
	}
	return net.ParseIP(addr)
}

type Map struct {
	// domains bypassed from adblock
	bypass    map[string]bool

	nBypass   int
	// domains to be blocked
	blocked   map[string]bool

	nBlocked  int
}

func newMap() *Map {
	m := new(Map)
	m.bypass = make(map[string]bool)
	m.blocked = make(map[string]bool)
	return m
}

func (m *Map) parse(r io.Reader){
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if i := bytes.Index(line, []byte{'#'}); i>=0 {
			line = line[:i]
		}
		f := bytes.Fields(line)
		if len(f) != 2 {
			continue
		}
		m.blocked[string(f[1])] = true
		m.nBlocked ++
	}
}

func openURL(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Warningf("Unable to read url %v", url)
		return nil, err
	}
	return resp.Body, nil
}