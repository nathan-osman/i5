package proxy

import (
	"github.com/nathan-osman/i5/geolocation"
	"github.com/nathan-osman/i5/notifier"
)

// Config provides the configuration for a proxy instance.
type Config struct {
	Addr        string
	Mountpoints []*Mountpoint
	Geolocation *geolocation.Geolocation
	Notifier    *notifier.Notifier
}
