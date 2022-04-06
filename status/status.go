package status

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/nathan-osman/go-herald"
	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/dbman"
	"github.com/nathan-osman/i5/dockmon"
	"github.com/nathan-osman/i5/ui"
	bolt "go.etcd.io/bbolt"
)

const (
	sessionName = "status"

	messageTypeRequest = "request"
)

// Status provides a set of endpoints that display status information.
type Status struct {
	Container *dockmon.Container
	conman    *conman.Conman
	dbman     *dbman.Manager
	db        *bolt.DB
	herald    *herald.Herald
	startup   int64
}

// New creates a new status container.
func New(cfg *Config) (*Status, error) {
	d, err := openDB(cfg.StorageDir)
	if err != nil {
		return nil, err
	}
	var (
		r = gin.Default()
		s = &Status{
			conman:  cfg.Conman,
			dbman:   cfg.Dbman,
			db:      d,
			herald:  herald.New(),
			startup: time.Now().Unix(),
		}
		store = cookie.NewStore([]byte(cfg.Key))
	)
	s.herald.Start()
	s.herald.SetCheckOrigin(func(r *http.Request) bool { return true })
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
	})
	r.Use(
		sessions.Sessions(sessionName, store),
		static.Serve("/", ui.EmbedFileSystem{FileSystem: http.FS(ui.Content)}),
	)
	if cfg.Debug {
		r.Use(
			cors.New(cors.Config{
				AllowOrigins:     []string{"http://localhost:3000"},
				AllowHeaders:     []string{"content-type"},
				AllowCredentials: true,
			}),
		)
	}
	auth := r.Group("/auth")
	{
		auth.POST("/login", s.authLogin)
		auth.POST("/logout", s.authLogout)
	}
	api := r.Group("/api")
	{
		api.Use(requireLogin)
		api.GET("/status", s.apiStatus)
		api.GET("/containers", s.apiContainers)
		api.GET("/ws", s.webSocket)
	}
	s.Container = &dockmon.Container{
		Domains:  []string{cfg.Domain},
		Insecure: cfg.Insecure,
		Handler:  r,
		Running:  true,
	}
	return s, nil
}

type messageRequest struct {
	RemoteAddr string `json:"remote_addr"`
	Method     string `json:"method"`
	Host       string `json:"host"`
	Path       string `json:"path"`
}

func (s *Status) BroadcastRequest(r *http.Request) {
	m, err := herald.NewMessage(messageTypeRequest, &messageRequest{
		RemoteAddr: r.RemoteAddr,
		Method:     r.Method,
		Host:       r.Host,
		Path:       r.URL.Path,
	})
	if err != nil {
		// TODO: determine if this is an appropriate response
		return
	}
	s.herald.Send(m, nil)
}

// Close shuts down the status server.
func (s *Status) Close() {
	s.herald.Close()
	s.db.Close()
}
