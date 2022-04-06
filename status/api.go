package status

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/nathan-osman/i5/conman"
)

var errWebSocketError = "unable to connect to WebSocket"

type apiStatusResponse struct {
	Startup int64 `json:"startup"`
}

func (s *Status) apiStatus(c *gin.Context) {
	c.JSON(http.StatusOK, &apiStatusResponse{
		Startup: s.startup,
	})
}

type byName []*conman.Info

func (n byName) Len() int           { return len(n) }
func (n byName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n byName) Less(i, j int) bool { return n[i].Name < n[j].Name }

func (s *Status) apiContainers(c *gin.Context) {
	containers := s.conman.Info()
	sort.Sort(byName(containers))
	c.JSON(http.StatusOK, containers)
}

func (s *Status) webSocket(c *gin.Context) {
	_, err := s.herald.AddClient(c.Writer, c.Request, nil)
	if err != nil {
		failure(c, 500, errWebSocketError)
		return
	}
}
