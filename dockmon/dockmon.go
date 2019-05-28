package dockmon

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/nathan-osman/i5/util"
	"github.com/sirupsen/logrus"
)

const (
	actionStart = "start"
	actionDie   = "die"
)

var evtOptions = types.EventsOptions{
	Filters: filters.NewArgs(),
}

func init() {
	evtOptions.Filters.Add("event", actionStart)
	evtOptions.Filters.Add("event", actionDie)
}

// Dockmon manages a connection to the Docker daemon and delivers events as containers are started and stopped. A list of running containers is also kept.
type Dockmon struct {
	log            *logrus.Entry
	client         *client.Client
	conMap         util.StringMap
	conStartedChan chan *Container
	conStoppedChan chan *Container
	closeFunc      context.CancelFunc
	closedChan     chan bool
}

func (d *Dockmon) add(ctx context.Context, id string) {
	if s, err := NewContainerFromClient(ctx, d.client, id); err == nil {
		d.conMap.Insert(id, s)
		d.conStartedChan <- s
	}
}

func (d *Dockmon) sync(ctx context.Context) error {
	d.log.Info("performing sync...")
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}
	// Create a StringMap of the containers
	conMap := util.StringMap{}
	for _, con := range containers {
		conMap.Insert(con.ID, nil)
	}
	// Add the containers that we don't know were started
	for conID := range conMap.Difference(d.conMap) {
		d.add(ctx, conID)
	}
	// Remove the containers that we don't know were stopped
	for containerID, s := range d.conMap.Difference(conMap) {
		d.conMap.Remove(containerID)
		d.conStoppedChan <- s.(*Container)
	}
	return nil
}

func (d *Dockmon) loop(ctx context.Context) error {
	// Watch for container start / stop events
	msgChan, errChan := d.client.Events(ctx, evtOptions)
	// ...but first sync with the daemon's list of running containers
	if err := d.sync(ctx); err != nil {
		return err
	}
	for {
		select {
		case msg := <-msgChan:
			switch msg.Action {
			case actionStart:
				d.add(ctx, msg.ID)
			case actionDie:
				if s, ok := d.conMap[msg.ID]; ok {
					d.conMap.Remove(msg.ID)
					d.conStoppedChan <- s.(*Container)
				}
			}
		case err := <-errChan:
			return err
		}
	}
}

func (d *Dockmon) run(ctx context.Context) {
	defer close(d.closedChan)
	defer close(d.conStoppedChan)
	defer close(d.conStartedChan)
	defer d.log.Info("Docker monitor stopped")
	d.log.Info("Docker monitor started")
	for {
		err := d.loop(ctx)
		if err == context.Canceled {
			return
		}
		d.log.Error(err)
		d.log.Info("reconnecting in 30 seconds...")
		select {
		case <-time.After(30 * time.Second):
		case <-ctx.Done():
			return
		}
	}
}

// New creates a new Dockmon instance and begins the event loop.
func New(cfg *Config) (*Dockmon, error) {
	c, err := client.NewClient(cfg.Host, "1.24", nil, nil)
	if err != nil {
		return nil, err
	}
	var (
		ctx, cancelFunc = context.WithCancel(context.Background())
		d               = &Dockmon{
			log:            logrus.WithField("context", "dockmon"),
			client:         c,
			conMap:         util.StringMap{},
			conStartedChan: make(chan *Container),
			conStoppedChan: make(chan *Container),
			closeFunc:      cancelFunc,
			closedChan:     make(chan bool),
		}
	)
	go d.run(ctx)
	return d, nil
}

// Monitor returns channels that send when containers are started and stopped.
func (d *Dockmon) Monitor() (<-chan *Container, <-chan *Container) {
	return d.conStartedChan, d.conStoppedChan
}

// Close shuts down the connection to the Docker daemon.
func (d *Dockmon) Close() {
	d.closeFunc()
	<-d.closedChan
}
