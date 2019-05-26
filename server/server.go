package server

import (
	"net"
	"net/http"
	"sync"

	"github.com/nathan-osman/i5/service"
	"github.com/sirupsen/logrus"
)

// Server listens for incoming connections and routes them accordingly.
type Server struct {
	mutex      sync.RWMutex
	log        *logrus.Entry
	listener   net.Listener
	domainMap  map[string]*service.Service
	closedChan chan bool
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	//...
}

func (s *Server) serve() {
	defer close(s.closedChan)
	defer s.log.Info("server stopped")
	s.log.Info("server started")
	server := http.Server{
		Handler: http.HandlerFunc(s.handle),
	}
	for {
		if err := server.Serve(s.listener); err == http.ErrServerClosed {
			return
		} else {
			s.log.Error(err)
		}
	}
}

// New creates a new server from the provided configuration.
func New(cfg *Config) (*Server, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	s := &Server{
		log:        logrus.WithField("context", "server"),
		listener:   l,
		domainMap:  map[string]*service.Service{},
		closedChan: make(chan bool),
	}
	go s.serve()
	//...
	return s, nil
}

// Close shuts down the server.
func (s *Server) Close() {
	_ = s.listener.Close()
	<-s.closedChan
}
