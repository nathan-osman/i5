package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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

func parsePort(addr string) (int, error) {
	if _, port, err := net.SplitHostPort(addr); err == nil {
		return net.LookupPort("tcp", port)
	} else {
		return 0, err
	}
}

func (s *Server) lookup(name string) (*dockmon.Container, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if v, ok := s.domainMap[strings.Split(name, ":")[0]]; ok {
		return v.(*dockmon.Container), nil
	} else {
		return nil, errInvalidDomain
	}
}

func (s *Server) decide(name string) error {
	_, err := s.lookup(name)
	return err
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request, con *dockmon.Container) {
	(&httputil.ReverseProxy{
		Director: func(inReq *http.Request) {
			inReq.Host = r.Host
			inReq.URL = &url.URL{
				Scheme: "http",
				Host:   con.Addr,
			}
		},
	}).ServeHTTP(w, r)
}

func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	if con, err := s.lookup(r.Host); err == nil {
		if con.Insecure {
			s.handle(w, r, con)
		} else {
			http.Redirect(
				w, r,
				fmt.Sprintf("https://%s%s", r.Host, r.URL.Path),
				http.StatusMovedPermanently,
			)
		}
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func (s *Server) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	if con, err := s.lookup(r.Host); err == nil {
		s.handle(w, r, con)
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
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
	// Set certmagic defaults
	s.cfg = certmagic.NewDefault()
	s.cfg.Agreed = true
	s.cfg.Email = cfg.Email
	s.cfg.OnDemand = &certmagic.OnDemandConfig{
		DecisionFunc: s.decide,
	}
	if cfg.Debug {
		s.cfg.CA = certmagic.LetsEncryptStagingCA
	}
	if port, err := parsePort(cfg.HTTPAddr); err == nil {
		s.cfg.AltHTTPPort = port
	} else {
		return nil, err
	}
	if len(cfg.StorageDir) != 0 {
		s.cfg.Storage = &certmagic.FileStorage{Path: cfg.StorageDir}
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
