package status

import (
	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/dbman"
)

// Config provides the configuration for the internal status server.
type Config struct {
	// Debug indicates that debugging is enabled.
	Debug bool
	// Domain indicates the domain that should be used for the container.
	Domain string
	// Insecure allows insecure connections to the server.
	Insecure bool
	// StorageDir indicates where the SQLite database should be stored
	StorageDir string
	// Conman is a pointer to a Conman instance.
	Conman *conman.Conman
	// Dbman is a pointer to a db.Manager instance.
	Dbman *dbman.Manager
}
