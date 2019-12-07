package status

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nathan-osman/i5/assets"
	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/db"
	"github.com/nathan-osman/i5/dockmon"
)

// Status provides a set of endpoints that display status information.
type Status struct {
	conman *conman.Conman
	dbman  *db.Manager
	router *mux.Router
}

// New creates a new status container.
func New(cfg *Config) *dockmon.Container {
	s := &Status{
		dbman:  cfg.Dbman,
		conman: cfg.Conman,
		router: mux.NewRouter(),
	}
	s.router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(assets.Assets)),
	)
	return &dockmon.Container{
		Domains:  []string{cfg.Domain},
		Insecure: cfg.Insecure,
		Handler:  s.router,
		Running:  true,
	}
}
