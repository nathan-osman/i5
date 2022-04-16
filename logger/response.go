package logger

import (
	"mime"
	"net"
	"net/http"

	"github.com/nathan-osman/go-herald"
)

const messageResponseType = "response"

type messageResponse struct {
	RemoteAddr    string `json:"remote_addr"`
	CountryCode   string `json:"country_code"`
	CountryName   string `json:"country_name"`
	Method        string `json:"method"`
	Host          string `json:"host"`
	Path          string `json:"path"`
	StatusCode    int    `json:"status_code"`
	Status        string `json:"status"`
	ContentLength string `json:"content_length"`
	ContentType   string `json:"content_type"`
}

// LogResponse logs an HTTP request and response.
func (l *Logger) LogResponse(resp *http.Response) {
	if l.herald != nil {
		host, _, err := net.SplitHostPort(resp.Request.RemoteAddr)
		if err != nil {
			host = resp.Request.RemoteAddr
		}
		var (
			countryCode string
			countryName string
		)
		if l.geolocator != nil {
			r, err := l.geolocator.Geolocate(host)
			if err == nil {
				countryCode = r.CountryCode
				countryName = r.CountryName
			}
		}
		contentType := resp.Header.Get("Content-Type")
		mediatype, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			mediatype = contentType
		}
		m, err := herald.NewMessage(messageResponseType, &messageResponse{
			RemoteAddr:    host,
			CountryCode:   countryCode,
			CountryName:   countryName,
			Method:        resp.Request.Method,
			Host:          resp.Request.Host,
			Path:          resp.Request.URL.Path,
			StatusCode:    resp.StatusCode,
			Status:        resp.Status,
			ContentLength: resp.Header.Get("Content-Length"),
			ContentType:   mediatype,
		})
		if err != nil {
			// TODO: log that the message couldn't be created
			return
		}
		l.herald.Send(m, nil)
	}
}
