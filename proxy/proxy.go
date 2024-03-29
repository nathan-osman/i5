package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/nathan-osman/i5/util"
)

const ContextSecure = "secure"

// CallbackFn provides a callback function for logging request data.
type CallbackFn func(resp *http.Response)

// Mountpoint represents a static path and directory on-disk.
type Mountpoint struct {
	Path string
	Dir  string
}

// Proxy acts as a reverse proxy and static file server for a specific site.
type Proxy struct {
	addr       string
	router     *chi.Mux
	callbackFn CallbackFn
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
			if p.callbackFn != nil {
				p.callbackFn(resp)
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
		addr:   cfg.Addr,
		router: chi.NewRouter(),
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

// SetCallback specifies a callback for requests & responses.
func (p *Proxy) SetCallback(callbackFn CallbackFn) {
	p.callbackFn = callbackFn
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}
