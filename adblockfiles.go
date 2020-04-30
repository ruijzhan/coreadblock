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

func (c *CoreAdBlock) parseHostsURL(url string) error  {
	reader, err := openURL(url)
	if err != nil {
		return err
	}
	defer reader.Close()
	c.parseHosts(reader)
	return nil
}

func (c *CoreAdBlock) parseHosts(r io.Reader){
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
		c.BlockList[string(f[1])] = true
		c.Bloom.Add([]byte(string(f[1])))
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