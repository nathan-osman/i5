package status

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/container"
	"github.com/nathan-osman/i5/db"
	"github.com/nathan-osman/i5/dbman"
	"github.com/nathan-osman/i5/dockmon"
	"github.com/nathan-osman/i5/logger"
	"github.com/nathan-osman/i5/ui"
)

const sessionName = "status"

// Status provides a set of endpoints that display status information.
type Status struct {
	Container *container.Container
	conman    *conman.Conman
	dockmon   *dockmon.Dockmon
	dbman     *dbman.Manager
	logger    *logger.Logger
	db        *db.DB
	startup   int64
}

// New creates a new status container.
func New(cfg *Config) *Status {
	var (
		r = gin.Default()
		s = &Status{
			conman:  cfg.Conman,
			dockmon: cfg.Dockmon,
			dbman:   cfg.Dbman,
			logger:  cfg.Logger,
			db:      cfg.DB,
			startup: time.Now().Unix(),
		}
		store = cookie.NewStore([]byte(cfg.Key))
	)
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
	})
	r.Use(
		sessions.Sessions(sessionName, store),
		static.Serve("/", ui.EmbedFileSystem{FileSystem: http.FS(ui.Content)}),
	)
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
		api.POST("/containers/:id/state", s.apiContainersState)
		api.GET("/db/:db/databases", s.apiDbDatabases)
		api.DELETE("/db/:db/databases/:name", s.apiDbDatabasesDelete)
		api.GET("/ws", s.webSocket)
	}
	r.NoRoute(func(c *gin.Context) {
		c.Request.URL.Path = "/"
		r.HandleContext(c)
		c.Abort()
	})
	s.Container = &container.Container{
		ContainerData: container.ContainerData{
			Domains:  []string{cfg.Domain},
			Insecure: cfg.Insecure,
		},
		Handler: r,
	}
	return s
}

// Close shuts down the status server.
func (s *Status) Close() {
	s.db.Close()
}
