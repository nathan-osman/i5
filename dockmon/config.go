package dockmon

import (
	"github.com/nathan-osman/geolocator"
	"github.com/nathan-osman/i5/notifier"
)

// Config provides the configuration for the Docker monitor.
type Config struct {
	// Host provides the address to use for connecting to the Docker daemon.
	Host string
	// Provider is a pointer to a geolocation.Provider instance.
	Provider geolocator.Provider
	// Notifier is a pointer to a Notifier instance.
	Notifier *notifier.Notifier
}
