package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"sync"

	"github.com/mholt/certmagic"
	"github.com/nathan-osman/i5/dockmon"
	"github.com/nathan-osman/i5/service"
	"github.com/nathan-osman/i5/util"
	"github.com/sirupsen/logrus"
)

var errInvalidDomain = errors.New("invalid domain name")

// Server listens for incoming connections and routes them accordingly.
type Server struct {
	mutex       sync.RWMutex
	log         *logrus.Entry
	cfg         *certmagic.Config
	dockmon     *dockmon.Dockmon
	httpServer  *http.Server
	httpsServer *http.Server
	domainMap   util.StringMap
	closeChan   chan bool
	closedChan  chan bool
}

func (s *Server) decide(name string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if !s.domainMap.Has(name) {
		return errInvalidDomain
	}
	return nil
}

func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	//...
}

func (s *Server) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	//...
}

func (s *Server) add(svc *service.Service) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, domain := range svc.Domains {
		s.domainMap.Insert(domain, svc)
	}
}

func (s *Server) remove(svc *service.Service) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, domain := range svc.Domains {
		s.domainMap.Remove(domain)
	}
}

func (s *Server) run() {
	defer close(s.closedChan)
	defer s.log.Info("server stopped")
	s.log.Info("server started")
	svcStartedChan, svcStoppedChan := s.dockmon.Monitor()
	for {
		select {
		case svc := <-svcStartedChan:
			s.add(svc)
		case svc := <-svcStoppedChan:
			s.remove(svc)
		case <-s.closeChan:
			return
		}
	}
}

// New creates a new server from the provided configuration.
func New(cfg *Config) (*Server, error) {
	// Create the server
	var s = &Server{
		log:         logrus.WithField("context", "server"),
		dockmon:     cfg.Dockmon,
		httpServer:  &http.Server{},
		httpsServer: &http.Server{},
		domainMap:   util.StringMap{},
		closeChan:   make(chan bool),
		closedChan:  make(chan bool),
	}
	// Finish initializing the server
	s.cfg = certmagic.NewDefault()
	s.cfg.Agreed = true
	s.cfg.Email = cfg.Email
	s.cfg.OnDemand = &certmagic.OnDemandConfig{
		DecisionFunc: s.decide,
	}
	if cfg.Debug {
		s.cfg.CA = certmagic.LetsEncryptStagingCA
	}
	s.httpServer.Handler = http.HandlerFunc(s.handleHTTP)
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
	// Launch goroutines for watching and for each of the servers
	go s.run()
	go func() {
		if err := s.httpServer.Serve(httpLn); err != http.ErrServerClosed {
			s.log.Error(err)
		}
	}()
	go func() {
		if err := s.httpsServer.Serve(httpsLn); err != http.ErrServerClosed {
			s.log.Error(err)
		}
	}()
	return s, nil
}

// Close shuts down the server.
func (s *Server) Close() {
	close(s.closeChan)
	<-s.closedChan
	s.log.Info("waiting for connections...")
	s.httpServer.Shutdown(context.Background())
	s.httpsServer.Shutdown(context.Background())
}
