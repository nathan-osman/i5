package proxy

import (
	"github.com/nathan-osman/geolocator"
	"github.com/nathan-osman/i5/notifier"
)

// Config provides the configuration for a proxy instance.
type Config struct {
	Addr        string
	Mountpoints []*Mountpoint
	Provider    geolocator.Provider
	Notifier    *notifier.Notifier
}
