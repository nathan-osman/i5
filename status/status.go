package status

import (
	"net/http"
	"sort"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/nathan-osman/i5/assets"
	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/dockmon"
)

// Status provides a set of endpoints that display status information.
type Status struct {
	conman      *conman.Conman
	router      *mux.Router
	templateSet *pongo2.TemplateSet
}

type ByName []*conman.Info

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return n[i].Name < n[j].Name }

func (s *Status) index(w http.ResponseWriter, r *http.Request) {
	i := s.conman.Info()
	sort.Sort(ByName(i))
	s.render(w, r, "index.html", pongo2.Context{
		"info": i,
	})
}

// New creates a new status container.
func New(cfg *Config) *dockmon.Container {
	s := &Status{
		conman:      cfg.Conman,
		router:      mux.NewRouter(),
		templateSet: pongo2.NewSet("", &vfsgenLoader{}),
	}
	s.router.PathPrefix("/static").Handler(http.FileServer(assets.Assets))
	s.router.HandleFunc("/", s.index)
	return &dockmon.Container{
		Domains:  []string{cfg.Domain},
		Insecure: cfg.Insecure,
		Handler:  s.router,
		Running:  true,
	}
}
