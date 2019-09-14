package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

const (
	ContextSecure = "secure"
)

// Mountpoint represents a static path and directory on-disk.
type Mountpoint struct {
	Path string
	Dir  string
}

// Proxy acts as a reverse proxy and static file server for a specific site.
type Proxy struct {
	addr   string
	router *mux.Router
}

func (p *Proxy) applyBranding(header http.Header) {
	header.Set("X-Powered-By", "i5 - qms.li/i5")
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
			p.applyBranding(resp.Header)
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			renderErrorTemplate(w, r, http.StatusBadGateway)
		},
	}).ServeHTTP(w, r)
}

// TODO: FileServer() should use applyBranding()

// New creates and initializes a new proxy for a domain.
func New(cfg *Config) *Proxy {
	p := &Proxy{
		addr:   cfg.Addr,
		router: mux.NewRouter(),
	}
	p.router.HandleFunc("/", p.handle)
	for _, m := range cfg.Mountpoints {
		p.router.PathPrefix(m.Path).Handler(
			http.FileServer(http.Dir(m.Dir)),
		)
	}
	return p
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}
