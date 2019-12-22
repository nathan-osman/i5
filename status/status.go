package status

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/db"
	"github.com/nathan-osman/i5/dockmon"
	"github.com/nathan-osman/i5/ui"
)

// Status provides a set of endpoints that display status information.
type Status struct {
	conman *conman.Conman
	dbman  *db.Manager
	router *chi.Mux
}

// New creates a new status container.
func New(cfg *Config) *dockmon.Container {
	s := &Status{
		dbman:  cfg.Dbman,
		conman: cfg.Conman,
		router: chi.NewRouter(),
	}
	s.router.Handle("/*", http.FileServer(ui.Assets))
	return &dockmon.Container{
		Domains:  []string{cfg.Domain},
		Insecure: cfg.Insecure,
		Handler:  s.router,
		Running:  true,
	}
}
