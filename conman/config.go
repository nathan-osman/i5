package conman

import (
	"github.com/nathan-osman/i5/db"
	"github.com/nathan-osman/i5/dockmon"
)

// Config provides the configuration for the container manager.
type Config struct {
	// EventChan receives container events.
	EventChan <-chan *dockmon.Event
	// Dbman is a pointer to a db.Manager instance.
	Dbman *db.Manager
}
