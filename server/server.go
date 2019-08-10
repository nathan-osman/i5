package server

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/mholt/certmagic"
	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/dockmon"
	"github.com/sirupsen/logrus"
)

// Server listens for incoming connections and routes them accordingly.
type Server struct {
	log         *logrus.Entry
	cfg         *certmagic.Config
	conman      *conman.Conman
	httpServer  *http.Server
	httpsServer *http.Server
}

func (s *Server) decide(name string) error {
	_, err := s.conman.Lookup(name)
	return err
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request, con *dockmon.Container, secure bool) {
	(&httputil.ReverseProxy{
		Director: func(inReq *http.Request) {
			if secure {
				inReq.Header.Set("X-Forwarded-Proto", "https")
			}
			if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
				inReq.Header.Set("X-Real-IP", host)
			}
			inReq.Host = r.Host
			inReq.URL = &url.URL{
				Scheme:   "http",
				Host:     con.Addr,
				Path:     r.URL.Path,
				RawQuery: r.URL.RawQuery,
			}
		},
		ModifyResponse: func(resp *http.Response) error {
			resp.Header.Set("X-Powered-By", "i5 - qms.li/i5")
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			renderErrorTemplate(w, r, http.StatusBadGateway)
		},
	}).ServeHTTP(w, r)
}

func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	if con, err := s.conman.Lookup(r.Host); err == nil {
		if con.Insecure {
			s.handle(w, r, con, false)
		} else {
			http.Redirect(
				w, r,
				(&url.URL{
					Scheme:   "https",
					Host:     r.Host,
					Path:     r.URL.Path,
					RawQuery: r.URL.RawQuery,
				}).String(),
				http.StatusMovedPermanently,
			)
		}
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func (s *Server) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	if con, err := s.conman.Lookup(r.Host); err == nil {
		s.handle(w, r, con, true)
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

// New creates a new server from the provided configuration.
func New(cfg *Config) (*Server, error) {

	// Create the server and certmagic config
	var (
		s = &Server{
			log:         logrus.WithField("context", "server"),
			conman:      cfg.Conman,
			httpServer:  &http.Server{},
			httpsServer: &http.Server{},
		}
		cmCfg = certmagic.Config{
			Agreed: true,
			Email:  cfg.Email,
			OnDemand: &certmagic.OnDemandConfig{
				DecisionFunc: s.decide,
			},
		}
	)

	// Finish initializing the config (the fields cannot be set inline)
	if cfg.Debug {
		cmCfg.CA = certmagic.LetsEncryptStagingCA
	}
	if len(cfg.StorageDir) != 0 {
		cmCfg.Storage = &certmagic.FileStorage{Path: cfg.StorageDir}
	}
	s.cfg = certmagic.New(
		certmagic.NewCache(
			certmagic.CacheOptions{
				GetConfigForCert: func(certmagic.Certificate) (certmagic.Config, error) {
					return certmagic.Default, nil
				},
			},
		),
		cmCfg,
	)
	s.httpServer.Handler = s.cfg.HTTPChallengeHandler(http.HandlerFunc(s.handleHTTP))
	s.httpsServer.Handler = http.HandlerFunc(s.handleHTTPS)

	// Create the HTTP listener
	httpLn, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		return nil, err
	}

	// Create the HTTPS listener
	httpsLn, err := tls.Listen("tcp", cfg.HTTPSAddr, s.cfg.TLSConfig())
	if err != nil {
		return nil, err
	}

	// Launch goroutines for each of the servers
	go func() {
		s.log.Info("listening for HTTP connections...")
		if err := s.httpServer.Serve(httpLn); err != http.ErrServerClosed {
			s.log.Errorf("HTTP: %s", err)
		}
	}()
	go func() {
		s.log.Info("listening for HTTPS connections...")
		if err := s.httpsServer.Serve(httpsLn); err != http.ErrServerClosed {
			s.log.Errorf("HTTPS: %s", err)
		}
	}()

	return s, nil
}

// Close shuts down the server.
func (s *Server) Close() {
	s.log.Info("shutting down server...")
	s.httpServer.Shutdown(context.Background())
	s.httpsServer.Shutdown(context.Background())
}
