package status

import (
	"github.com/nathan-osman/i5/conman"
)

// Config provides the configuration for the internal status server.
type Config struct {
	// Domain indicates the domain that should be used for the container.
	Domain string
	// Insecure allows insecure connections to the server.
	Insecure bool
	// Conman is a pointer to a Conman instance.
	Conman *conman.Conman
}
