package util

import (
	"net"
	"strings"
)

// ParseHost parses an address of the form "host:port" and returns the hist portion.
func ParseHost(addr string) string {
	return strings.SplitN(addr, ":", 2)[0]
}

// ParsePort parses an address of the form "host:port" and returns the port portion. Special service names (such as "http") are recognized.
func ParsePort(addr string, def int) (int, error) {
	if parts := strings.SplitN(addr, ":", 2); len(parts) == 2 {
		return net.LookupPort("tcp", parts[1])
	} else {
		return def, nil
	}
}
