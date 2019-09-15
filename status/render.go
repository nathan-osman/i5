package status

import (
	"net/http"
	"strconv"

	"github.com/flosch/pongo2"
)

func (s *Status) render(w http.ResponseWriter, r *http.Request, templateName string, ctx pongo2.Context) {
	err := func() error {
		t, err := s.templateSet.FromFile(templateName)
		if err != nil {
			return err
		}
		b, err := t.ExecuteBytes(ctx)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(b)))
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return nil
	}()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
