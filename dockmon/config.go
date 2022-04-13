package dockmon

import (
	"github.com/nathan-osman/i5/geolocation"
	"github.com/nathan-osman/i5/notifier"
)

// Config provides the configuration for the Docker monitor.
type Config struct {
	// Host provides the address to use for connecting to the Docker daemon.
	Host string
	// Geolocation is a pointer to a Geolocation instance.
	Geolocation *geolocation.Geolocation
	// Notifier is a pointer to a Notifier instance.
	Notifier *notifier.Notifier
}
