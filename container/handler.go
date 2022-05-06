package container

import (
	"net/http"

	"github.com/nathan-osman/i5/util"
)

type disabledHandler struct {
	Message string
}

func (d *disabledHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.RenderError(w, r, http.StatusInternalServerError, d.Message)
}
