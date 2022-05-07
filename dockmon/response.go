package dockmon

import (
	"mime"
	"net"
	"net/http"
)

const messageTypeResponse = "response"

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

func (d *Dockmon) responseCallback(resp *http.Response) {
	host, _, err := net.SplitHostPort(resp.Request.RemoteAddr)
	if err != nil {
		host = resp.Request.RemoteAddr
	}
	var (
		countryCode string
		countryName string
	)
	r, err := d.logger.Geolocate(host)
	if err == nil {
		countryCode = r.CountryCode
		countryName = r.CountryName
	}
	contentType := resp.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		mediatype = contentType
	}
	d.logger.Log(messageTypeResponse, &messageResponse{
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
}
