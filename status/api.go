package status

import (
	"net/http"

	"github.com/go-chi/render"
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

type containerResponse struct {
	Name    string `json:"name"`
	Running bool   `json:"running"`
}

func (c *containerResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Status) getContainers(w http.ResponseWriter, r *http.Request) {
	containers := []render.Renderer{}
	for _, c := range s.conman.Info() {
		containers = append(containers, &containerResponse{
			Name:    c.Name,
			Running: c.Running,
		})
	}
	render.RenderList(w, r, containers)
}
