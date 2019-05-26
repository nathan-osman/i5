package server

import (
	"github.com/nathan-osman/i5/dockmon"
)

// Config provides the configuration for the i5 server.
type Config struct {
	// Debug enables debug mode.
	Debug bool
	// Email provides the address to provide with challenges.
	Email string
	// HTTPAddr indicates the address to use for HTTP connections.
	HTTPAddr string
	// HTTPSAddr indicates the address to use for HTTPS connections.
	HTTPSAddr string
	// Dockmon is a pointer to a Dockmon instance.
	Dockmon *dockmon.Dockmon
}
