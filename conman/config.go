package conman

import (
	"github.com/nathan-osman/i5/dockmon"
)

// Config provides the configuration for the container manager.
type Config struct {
	// ConStartedChan receives container start events.
	ConStartedChan <-chan *dockmon.Container
	// ConStoppedChan receives container stop events.
	ConStoppedChan <-chan *dockmon.Container
}
