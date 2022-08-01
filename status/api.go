package status

import (
	"context"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/nathan-osman/i5/container"
)

const (
	containerActionStart = "start"
	containerActionStop  = "stop"

	errInvalidAction = "invalid action specified"
)

type apiStatusDatabaseResponse struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	Version string `json:"version"`
}

type apiStatusResponse struct {
	Startup   int64                        `json:"startup"`
	Databases []*apiStatusDatabaseResponse `json:"databases"`
}

func (s *Status) apiStatus(c *gin.Context) {
	var (
		l         = s.dbman.List()
		databases = []*apiStatusDatabaseResponse{}
	)
	for _, d := range l {
		databases = append(
			databases,
			&apiStatusDatabaseResponse{
				Name:    d.Name(),
				Title:   d.Title(),
				Version: d.Version(),
			},
		)
	}
	c.JSON(http.StatusOK, &apiStatusResponse{
		Startup:   s.startup,
		Databases: databases,
	})
}

type byName []container.ContainerData

func (n byName) Len() int           { return len(n) }
func (n byName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n byName) Less(i, j int) bool { return n[i].Name < n[j].Name }

func (s *Status) apiContainers(c *gin.Context) {
	containers := s.conman.Info()
	sort.Sort(byName(containers))
	c.JSON(http.StatusOK, containers)
}

type apiContainerStateParams struct {
	Action string `json:"action"`
}

func (s *Status) apiContainersState(c *gin.Context) {
	var (
		id     = c.Param("id")
		params = &apiContainerStateParams{}
	)
	if err := c.ShouldBindJSON(params); err != nil {
		failure(c, http.StatusBadRequest, err.Error())
		return
	}
	var fn func(context.Context, string) error
	switch params.Action {
	case containerActionStart:
		fn = s.dockmon.StartContainer
	case containerActionStop:
		fn = s.dockmon.StopContainer
	default:
		failure(c, http.StatusBadRequest, errInvalidAction)
		return
	}
	if err := fn(c.Request.Context(), id); err != nil {
		failure(c, http.StatusInternalServerError, err.Error())
		return
	}
	success(c)
}

func (s *Status) webSocket(c *gin.Context) {
	if err := s.logger.AddClient(c.Writer, c.Request); err != nil {
		// TODO: not sure if we can return a response here since we don't know
		// what state the connection is in, so just do nothing for the moment
		return
	}
}
