package server

import (
	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/notifier"
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
	// StorageDir indicates where certificates should be stored
	StorageDir string
	// Conman is a pointer to a Conman instance.
	Conman *conman.Conman
	// Notifier is a pointer to a Notifier instance.
	Notifier *notifier.Notifier
}
