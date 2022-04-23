package logger

import (
	"net/http"

	"github.com/nathan-osman/geolocator"
	"github.com/nathan-osman/geolocator/ip2location"
	"github.com/nathan-osman/go-herald"
)

// Logger acts as a centralized component for ingesting data from the other
// components and processing it.
type Logger struct {
	geolocator geolocator.Provider
	herald     *herald.Herald
}

func createProvider(dbType, path string) (geolocator.Provider, error) {
	switch dbType {
	case "ip2location":
		return ip2location.New(path)
	}
	return nil, nil
}

// New creates a new logger instance.
func New(cfg *Config) (*Logger, error) {
	p, err := createProvider(cfg.GeolocationDBType, cfg.GeolocationDBPath)
	if err != nil {
		return nil, err
	}
	var h *herald.Herald
	if cfg.Status {
		h = herald.New()
		h.SetCheckOrigin(func(r *http.Request) bool { return true })
	}
	return &Logger{
		geolocator: p,
		herald:     h,
	}, nil
}

// AddClient upgrades the provided HTTP request to websocket.
func (l *Logger) AddClient(w http.ResponseWriter, r *http.Request) error {
	_, err := l.herald.AddClient(w, r, nil)
	return err
}

// Close shuts down the logger.
func (l *Logger) Close() {
	if l.herald != nil {
		l.herald.Close()
	}
	if l.geolocator != nil {
		l.geolocator.Close()
	}
}
