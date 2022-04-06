package server

import (
	"net/http"

	"github.com/nathan-osman/i5/conman"
)

// HookFn is a function called when a new request comes in.
type HookFn func(*http.Request)

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
	// Hook (if not nil) is invoked once per request.
	Hook HookFn
	// Conman is a pointer to a Conman instance.
	Conman *conman.Conman
}
