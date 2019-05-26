package server

import (
	"github.com/nathan-osman/i5/dockmon"
)

// Config provides the configuration for the i5 server.
type Config struct {
	// Addr indicates the address to be used for listening.
	Addr string
	// Dockmon is a pointer to a Dockmon instance.
	Dockmon *dockmon.Dockmon
}
