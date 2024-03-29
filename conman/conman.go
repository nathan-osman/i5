package conman

import (
	"fmt"
	"sync"

	"github.com/nathan-osman/i5/container"
	"github.com/nathan-osman/i5/dbman"
	"github.com/nathan-osman/i5/dockmon"
	"github.com/nathan-osman/i5/util"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidDomain = fmt.Errorf("invalid domain name specified")
)

// Conman manages running containers based on Docker events.
type Conman struct {
	mutex      sync.RWMutex
	log        *logrus.Entry
	idMap      util.StringMap
	domainMap  util.StringMap
	eventChan  <-chan *dockmon.Event
	dbman      *dbman.Manager
	closeChan  chan bool
	closedChan chan bool
}

func (c *Conman) initDatabase(con *container.Container) error {
	d, err := c.dbman.Get(con.Database.Driver)
	if err != nil {
		return err
	}
	if err := d.CreateUser(
		con.Database.User,
		con.Database.Password,
	); err != nil {
		return err
	}
	if err := d.CreateDatabase(
		con.Database.Name,
		con.Database.User,
	); err != nil {
		return err
	}
	return nil
}

func (c *Conman) run() {
	defer close(c.closedChan)
	defer c.log.Info("service manager stopped")
	c.log.Info("service manager started")
	for {
		select {
		case e := <-c.eventChan:
			switch e.Action {
			case dockmon.Create:
				c.Add(e.Container)
			case dockmon.Destroy:
				c.Remove(e.Container)
			case dockmon.Start:
				c.ToggleState(e.Container, true)
			case dockmon.Die:
				c.ToggleState(e.Container, false)
			}
		case <-c.closeChan:
			return
		}
	}
}

// New creates a new container manager.
func New(cfg *Config) *Conman {
	c := &Conman{
		log:        logrus.WithField("context", "conman"),
		idMap:      util.StringMap{},
		domainMap:  util.StringMap{},
		eventChan:  cfg.EventChan,
		dbman:      cfg.Dbman,
		closeChan:  make(chan bool),
		closedChan: make(chan bool),
	}
	go c.run()
	return c
}

func (c *Conman) Add(con *container.Container) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.idMap.Insert(con.ID, con)
	for _, domain := range con.Domains {
		c.log.Debugf("created container for %s", domain)
		c.domainMap.Insert(domain, con)
	}
	if con.Database != nil {
		if err := c.initDatabase(con); err != nil {
			c.log.Error(err)
		}
	}
}

func (c *Conman) Remove(con *container.Container) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.idMap.Remove(con.ID)
	for _, domain := range con.Domains {
		c.log.Debugf("removed container for %s", domain)
		c.domainMap.Remove(domain)
	}
}

func (c *Conman) ToggleState(con *container.Container, running bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.idMap.Insert(con.ID, con)
	for _, domain := range con.Domains {
		c.log.Debugf("toggled %s running to %t", domain, running)
		c.domainMap.Insert(domain, con)
	}
}

// Lookup attempts to retrieve the container for the provided domain name.
func (c *Conman) Lookup(name string) (*container.Container, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if v, ok := c.domainMap[name]; ok {
		return v.(*container.Container), nil
	} else {
		return nil, ErrInvalidDomain
	}
}

// Close shuts down the container manager.
func (c *Conman) Close() {
	close(c.closeChan)
	<-c.closedChan
}
