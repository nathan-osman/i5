package util

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/nathan-osman/i5/assets"
)

var errorTemplate *template.Template

func init() {
	f, err := assets.Assets.Open("templates/error.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	errorTemplate = template.Must(template.New("").Parse(string(b)))
}

// RenderError renders the error template with the specified error.
func RenderError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	buff := &bytes.Buffer{}
	errorTemplate.Execute(buff, message)
	b := buff.Bytes()
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(b)
}
