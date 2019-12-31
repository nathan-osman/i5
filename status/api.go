package status

import (
	"net/http"
	"sort"

	"github.com/go-chi/render"
	"github.com/nathan-osman/i5/conman"
)

type statusResponse struct {
	Startup int64 `json:"startup"`
}

func (s *statusResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Status) getStatus(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, &statusResponse{
		Startup: s.startup,
	})
}

type infoResponse struct{ *conman.Info }

func (i *infoResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type byName []*conman.Info

func (n byName) Len() int           { return len(n) }
func (n byName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n byName) Less(i, j int) bool { return n[i].Name < n[j].Name }

func (s *Status) getContainers(w http.ResponseWriter, r *http.Request) {
	var (
		containers    = s.conman.Info()
		containerList = []render.Renderer{}
	)
	sort.Sort(byName(containers))
	for _, c := range containers {
		containerList = append(containerList, &infoResponse{c})
	}
	render.RenderList(w, r, containerList)
}
