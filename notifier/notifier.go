package notifier

import (
	"net/http"

	"github.com/nathan-osman/go-herald"
)

// Notifier provides a simple way for components to send status messages on a
// WebSocket.
type Notifier struct {
	*herald.Herald
}

// New creates a new
func New(cfg *Config) *Notifier {
	n := &Notifier{
		herald.New(),
	}
	n.Start()
	if cfg.Debug {
		n.SetCheckOrigin(func(r *http.Request) bool { return true })
	}
	return n
}
