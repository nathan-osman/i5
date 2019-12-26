package status

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/db"
	"github.com/nathan-osman/i5/dockmon"
	"github.com/nathan-osman/i5/ui"
)

// Status provides a set of endpoints that display status information.
type Status struct {
	conman  *conman.Conman
	dbman   *db.Manager
	startup int64
}

// New creates a new status container.
func New(cfg *Config) *dockmon.Container {
	s := &Status{
		dbman:   cfg.Dbman,
		conman:  cfg.Conman,
		startup: time.Now().Unix(),
	}
	router := chi.NewRouter()
	router.Mount("/", http.FileServer(ui.Assets))
	router.Route("/api", func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Get("/status", s.getStatus)
		r.Get("/containers", s.getContainers)
	})
	return &dockmon.Container{
		Domains:  []string{cfg.Domain},
		Insecure: cfg.Insecure,
		Handler:  router,
		Running:  true,
	}
}
