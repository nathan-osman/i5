package proxy

import (
	"mime"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/nathan-osman/geolocator"
	"github.com/nathan-osman/go-herald"
	"github.com/nathan-osman/i5/notifier"
	"github.com/nathan-osman/i5/util"
)

const (
	ContextSecure = "secure"

	messageTypeRequest = "request"
)

// Mountpoint represents a static path and directory on-disk.
type Mountpoint struct {
	Path string
	Dir  string
}

// Proxy acts as a reverse proxy and static file server for a specific site.
type Proxy struct {
	addr     string
	router   *chi.Mux
	provider geolocator.Provider
	notifier *notifier.Notifier
}

func applyBranding(header http.Header) {
	header.Set("X-Powered-By", "i5 - qms.li/i5")
}

func brandingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		applyBranding(w.Header())
		h.ServeHTTP(w, r)
	})
}

type messageRequest struct {
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

func (p *Proxy) sendRequest(resp *http.Response) {
	host, _, err := net.SplitHostPort(resp.Request.RemoteAddr)
	if err != nil {
		host = resp.Request.RemoteAddr
	}
	var (
		countryCode string
		countryName string
	)
	if p.provider != nil {
		r, err := p.provider.Geolocate(host)
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
	m, err := herald.NewMessage(messageTypeRequest, &messageRequest{
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
		return
	}
	p.notifier.Send(m, nil)
}

func (p *Proxy) handle(w http.ResponseWriter, r *http.Request) {
	(&httputil.ReverseProxy{
		Director: func(inReq *http.Request) {
			if r.Context().Value(ContextSecure) != nil {
				inReq.Header.Set("X-Forwarded-Proto", "https")
			}
			if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
				inReq.Header.Set("X-Real-IP", host)
			}
			inReq.Host = r.Host
			inReq.URL = &url.URL{
				Scheme:   "http",
				Host:     p.addr,
				Path:     r.URL.Path,
				RawQuery: r.URL.RawQuery,
			}
		},
		ModifyResponse: func(resp *http.Response) error {
			if p.notifier != nil {
				p.sendRequest(resp)
			}
			applyBranding(resp.Header)
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			util.RenderError(w, r, http.StatusBadGateway, err.Error())
		},
	}).ServeHTTP(w, r)
}

// New creates and initializes a new proxy for a domain.
func New(cfg *Config) *Proxy {
	p := &Proxy{
		addr:     cfg.Addr,
		router:   chi.NewRouter(),
		provider: cfg.Provider,
		notifier: cfg.Notifier,
	}
	for _, m := range cfg.Mountpoints {
		p.router.Handle(
			m.Path+"*",
			brandingHandler(
				http.StripPrefix(
					m.Path,
					http.FileServer(http.Dir(m.Dir)),
				),
			),
		)
	}
	if p.addr != "" {
		p.router.HandleFunc("/*", p.handle)
	}
	return p
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}
