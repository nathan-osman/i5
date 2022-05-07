package logger

import (
	"errors"
	"net/http"

	"github.com/nathan-osman/geolocator"
	"github.com/nathan-osman/geolocator/ip2location"
	"github.com/nathan-osman/go-herald"
)

var (
	errNoGeolocator = errors.New("no geolocation database provided")
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
	h := herald.New()
	h.SetCheckOrigin(func(r *http.Request) bool { return true })
	h.Start()
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

func (l *Logger) Geolocate(host string) (*geolocator.Response, error) {
	if l.geolocator != nil {
		return l.geolocator.Geolocate(host)
	}
	return nil, errNoGeolocator
}

func (l *Logger) Log(action string, data interface{}) {
	m, err := herald.NewMessage(action, data)
	if err != nil {
		// TODO: something other than panic here
		panic(err)
	}
	l.herald.Send(m, nil)
}

// Close shuts down the logger.
func (l *Logger) Close() {
	l.herald.Close()
	if l.geolocator != nil {
		l.geolocator.Close()
	}
}
