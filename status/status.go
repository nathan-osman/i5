package status

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/dbman"
	"github.com/nathan-osman/i5/dockmon"
	"github.com/nathan-osman/i5/status/db"
	"github.com/nathan-osman/i5/ui"
)

// Status provides a set of endpoints that display status information.
type Status struct {
	Container *dockmon.Container
	conman    *conman.Conman
	dbman     *dbman.Manager
	conn      *db.Conn
	startup   int64
}

// New creates a new status container.
func New(cfg *Config) (*Status, error) {
	d, err := db.New(cfg.StorageDir)
	if err != nil {
		return nil, err
	}
	s := &Status{
		conman:  cfg.Conman,
		dbman:   cfg.Dbman,
		conn:    d,
		startup: time.Now().Unix(),
	}
	router := chi.NewRouter()
	router.Mount("/", http.FileServer(http.FS(ui.Content)))
	router.Route("/api", func(r chi.Router) {
		r.Get("/status", s.getStatus)
		r.Get("/containers", s.getContainers)
	})
	s.Container = &dockmon.Container{
		Domains:  []string{cfg.Domain},
		Insecure: cfg.Insecure,
		Handler:  router,
		Running:  true,
	}
	return s, nil
}

// Close shuts down the status server.
func (s *Status) Close() {
	s.conn.Close()
}
