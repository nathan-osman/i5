package dockmon

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/nathan-osman/i5/service"
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
	svcMap         util.StringMap
	svcStartedChan chan *service.Service
	svcStoppedChan chan *service.Service
	closeFunc      context.CancelFunc
	closedChan     chan bool
}

func (d *Dockmon) add(ctx context.Context, id string) error {
	s, err := service.NewFromContainer(ctx, d.client, id)
	if err != nil {
		return err
	}
	d.svcMap.Insert(id, s)
	d.svcStartedChan <- s
	return nil
}

func (d *Dockmon) sync(ctx context.Context) error {
	d.log.Info("performing sync...")
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}
	// Create a StringMap of the containers
	containerMap := util.StringMap{}
	for _, container := range containers {
		containerMap.Insert(container.ID, nil)
	}
	// Add the containers that we don't know were started
	for containerID := range containerMap.Difference(d.svcMap) {
		if err := d.add(ctx, containerID); err != nil {
			return err
		}
	}
	// Remove the containers that we don't know were stopped
	for containerID, s := range d.svcMap.Difference(containerMap) {
		d.svcMap.Remove(containerID)
		d.svcStoppedChan <- s.(*service.Service)
	}
	return nil
}

func (d *Dockmon) loop(ctx context.Context) error {
	// Sync with the daemon's list of running containers
	if err := d.sync(ctx); err != nil {
		return err
	}
	// Watch for container start / stop events
	msgChan, errChan := d.client.Events(ctx, evtOptions)
	for {
		select {
		case msg := <-msgChan:
			switch msg.Action {
			case actionStart:
				if err := d.add(ctx, msg.ID); err != nil {
					return err
				}
			case actionDie:
				if s, ok := d.svcMap[msg.ID]; ok {
					d.svcMap.Remove(msg.ID)
					d.svcStoppedChan <- s.(*service.Service)
				}
			}
		case err := <-errChan:
			return err
		}
	}
}

func (d *Dockmon) run(ctx context.Context) {
	defer close(d.closedChan)
	defer close(d.svcStoppedChan)
	defer close(d.svcStartedChan)
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
			svcMap:         util.StringMap{},
			svcStartedChan: make(chan *service.Service),
			svcStoppedChan: make(chan *service.Service),
			closeFunc:      cancelFunc,
			closedChan:     make(chan bool),
		}
	)
	go d.run(ctx)
	return d, nil
}

// Monitor returns channels that send when services are started and stopped.
func (d *Dockmon) Monitor() (<-chan *service.Service, <-chan *service.Service) {
	return d.svcStartedChan, d.svcStoppedChan
}

// Close shuts down the connection to the Docker daemon.
func (d *Dockmon) Close() {
	d.closeFunc()
	<-d.closedChan
}
