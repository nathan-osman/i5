package server

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"path"

	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/proxy"
	"github.com/nathan-osman/i5/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

const (
	errInvalidDomainName = "invalid domain name specified"

	letsEncryptStagingURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
)

// Server listens for incoming connections and routes them accordingly.
type Server struct {
	log         *logrus.Entry
	conman      *conman.Conman
	httpServer  *http.Server
	httpsServer *http.Server
}

func (s *Server) decide(ctx context.Context, host string) error {
	_, err := s.conman.Lookup(host)
	return err
}

func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	if con, err := s.conman.Lookup(r.Host); err == nil {
		if con.Insecure {
			con.Handler.ServeHTTP(w, r)
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
		con.Handler.ServeHTTP(w, r.WithContext(
			context.WithValue(r.Context(), proxy.ContextSecure, true),
		))
	} else {
		util.RenderError(w, r, http.StatusBadRequest, errInvalidDomainName)
	}
}

// New creates a new server from the provided configuration.
func New(cfg *Config) (*Server, error) {

	// Create the server and autocert config
	var (
		s = &Server{
			log:         logrus.WithField("context", "server"),
			conman:      cfg.Conman,
			httpServer:  &http.Server{},
			httpsServer: &http.Server{},
		}
		manager = autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache(path.Join(cfg.StorageDir, "certificates")),
			HostPolicy: s.decide,
			Email:      cfg.Email,
		}
	)

	// If debug mode is enabled, use the staging URL
	if cfg.Debug {
		manager.Client = &acme.Client{
			DirectoryURL: letsEncryptStagingURL,
		}
	}

	s.httpServer.Handler = manager.HTTPHandler(http.HandlerFunc(s.handleHTTP))
	s.httpsServer.Handler = http.HandlerFunc(s.handleHTTPS)

	// Create the HTTP listener
	httpListener, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		return nil, err
	}

	// Create the HTTPS listener
	httpsListener, err := tls.Listen("tcp", cfg.HTTPSAddr, manager.TLSConfig())
	if err != nil {
		return nil, err
	}

	// Launch goroutines for each of the servers
	go func() {
		s.log.Info("listening for HTTP connections...")
		if err := s.httpServer.Serve(httpListener); err != http.ErrServerClosed {
			s.log.Errorf("HTTP: %s", err)
		}
	}()
	go func() {
		s.log.Info("listening for HTTPS connections...")
		if err := s.httpsServer.Serve(httpsListener); err != http.ErrServerClosed {
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
