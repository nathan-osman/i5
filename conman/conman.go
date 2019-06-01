package conman

import (
	"fmt"
	"sync"

	"github.com/nathan-osman/i5/dockmon"
	"github.com/nathan-osman/i5/util"
	"github.com/sirupsen/logrus"
)

// ConMan manages running containers based on Docker events.
type ConMan struct {
	mutex          sync.RWMutex
	log            *logrus.Entry
	domainMap      util.StringMap
	conStartedChan <-chan *dockmon.Container
	conStoppedChan <-chan *dockmon.Container
	closeChan      chan bool
	closedChan     chan bool
}

func (c *ConMan) add(con *dockmon.Container) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, domain := range con.Domains {
		c.log.Debugf("added %s", domain)
		c.domainMap.Insert(domain, con)
	}
}

func (c *ConMan) remove(con *dockmon.Container) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, domain := range con.Domains {
		c.log.Debugf("removed %s", domain)
		c.domainMap.Remove(domain)
	}
}

func (c *ConMan) run() {
	defer close(c.closedChan)
	defer c.log.Info("service manager stopped")
	c.log.Info("service manager started")
	for {
		select {
		case con := <-c.conStartedChan:
			c.add(con)
		case con := <-c.conStoppedChan:
			c.remove(con)
		case <-c.closeChan:
			return
		}
	}
}

// New creates a new container manager.
func New(cfg *Config) *ConMan {
	c := &ConMan{
		log:            logrus.WithField("context", "conman"),
		domainMap:      util.StringMap{},
		conStartedChan: cfg.ConStartedChan,
		conStoppedChan: cfg.ConStoppedChan,
		closeChan:      make(chan bool),
		closedChan:     make(chan bool),
	}
	go c.run()
	return c
}

// Lookup attempts to retrieve the container for the provided domain name.
func (c *ConMan) Lookup(name string) (*dockmon.Container, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if v, ok := c.domainMap[name]; ok {
		return v.(*dockmon.Container), nil
	} else {
		return nil, fmt.Errorf("invalid domain \"%s\"", name)
	}
}

// Close shuts down the container manager.
func (c *ConMan) Close() {
	close(c.closeChan)
	<-c.closedChan
}
