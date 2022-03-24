package status

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/dbman"
	"github.com/nathan-osman/i5/dockmon"
	"github.com/nathan-osman/i5/ui"
	bolt "go.etcd.io/bbolt"
)

// Status provides a set of endpoints that display status information.
type Status struct {
	Container *dockmon.Container
	conman    *conman.Conman
	dbman     *dbman.Manager
	db        *bolt.DB
	startup   int64
}

// New creates a new status container.
func New(cfg *Config) (*Status, error) {
	d, err := openDB(cfg.StorageDir)
	if err != nil {
		return nil, err
	}
	s := &Status{
		conman:  cfg.Conman,
		dbman:   cfg.Dbman,
		db:      d,
		startup: time.Now().Unix(),
	}
	router := chi.NewRouter()
	router.Mount("/", http.FileServer(http.FS(ui.Content)))
	router.Route("/api", func(r chi.Router) {
		if cfg.Debug {
			r.Use(cors.Handler(
				cors.Options{
					AllowedOrigins: []string{"http://localhost:3000"},
				},
			))
		}
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
	s.db.Close()
}
