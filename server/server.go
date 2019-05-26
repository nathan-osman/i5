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

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
}

func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	s.handle(w, r)
}

func (s *Server) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	s.handle(w, r)
}

func (s *Server) add(con *dockmon.Container) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, domain := range con.Domains {
		s.log.Debugf("added %s", domain)
		s.domainMap.Insert(domain, con)
	}
}

func (s *Server) remove(con *dockmon.Container) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, domain := range con.Domains {
		s.log.Debugf("removed %s", domain)
		s.domainMap.Remove(domain)
	}
}

func (s *Server) run() {
	defer close(s.closedChan)
	defer s.log.Info("server stopped")
	s.log.Info("server started")
	conStartedChan, conStoppedChan := s.dockmon.Monitor()
	for {
		select {
		case con := <-conStartedChan:
			s.add(con)
		case con := <-conStoppedChan:
			s.remove(con)
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
