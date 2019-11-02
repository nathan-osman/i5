package server

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"

	"github.com/mholt/certmagic"
	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/dockmon"
	"github.com/nathan-osman/i5/proxy"
	"github.com/nathan-osman/i5/util"
	"github.com/sirupsen/logrus"
)

const (
	errContainerNotRunning = "container serving this domain is not running"
	errInvalidDomainName   = "invalid domain name specified"
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

func (s *Server) handleRequest(con *dockmon.Container, w http.ResponseWriter, r *http.Request) {
	if con.Running {
		con.Handler.ServeHTTP(w, r)
	} else {
		util.RenderError(w, r, http.StatusInternalServerError, errContainerNotRunning)
	}
}

func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	if con, err := s.conman.Lookup(r.Host); err == nil {
		if con.Insecure {
			s.handleRequest(con, w, r)
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
		util.RenderError(w, r, http.StatusBadRequest, errInvalidDomainName)
	}
}

func (s *Server) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	if con, err := s.conman.Lookup(r.Host); err == nil {
		s.handleRequest(con, w, r.WithContext(
			context.WithValue(r.Context(), proxy.ContextSecure, true),
		))
	} else {
		util.RenderError(w, r, http.StatusBadRequest, errInvalidDomainName)
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
