package geolocation

import (
	"errors"
)

var errInvalidGeolocationDBType = errors.New("invalid geolocation database type")

type geolocationProvider interface {
	Country(string) string
	Close()
}

// Geolocation provides a unified interface for interacting with geolocation
// databases on disk.
type Geolocation struct {
	provider geolocationProvider
}

// New creates a new instance of Geolocation.
func New(cfg *Config) (*Geolocation, error) {
	var provider geolocationProvider
	switch cfg.DBType {
	case ip2locationType:
		p, err := newIp2locationProvider(cfg.DBPath)
		if err != nil {
			return nil, err
		}
		provider = p
	default:
		return nil, errInvalidGeolocationDBType
	}
	return &Geolocation{
		provider: provider,
	}, nil
}

// Country returns the 2 character ISO country code for the provided IP address
// if available, or an empty string if not.
func (g *Geolocation) Country(ip string) string {
	return g.provider.Country(ip)
}

// Close frees resources being used by the geolocation service.
func (g *Geolocation) Close() {
	g.provider.Close()
}
