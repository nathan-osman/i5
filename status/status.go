package status

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/dockmon"
)

// Status provides a set of endpoints that display status information.
type Status struct {
	conman *conman.Conman
	router *mux.Router
}

func (s *Status) index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

// New creates a new status container.
func New(cfg *Config) *dockmon.Container {
	s := &Status{
		conman: cfg.Conman,
		router: mux.NewRouter(),
	}
	s.router.HandleFunc("/", s.index)
	return &dockmon.Container{
		Domains: []string{cfg.Domain},
		Handler: s.router,
	}
}
