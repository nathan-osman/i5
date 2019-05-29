package util

import (
	"net"
	"strings"
)

// ParsePort parses an address of the form "host:port" and returns the port portion. Recognizable service names (such as "http") are recognized.
func ParsePort(addr string, def int) (int, error) {
	if parts := strings.SplitN(addr, ":", 2); len(parts) == 2 {
		return net.LookupPort("tcp", parts[1])
	} else {
		return def, nil
	}
}
