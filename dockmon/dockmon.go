package dockmon

import (
	"context"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/nathan-osman/i5/container"
	"github.com/nathan-osman/i5/logger"
	"github.com/nathan-osman/i5/util"
	"github.com/sirupsen/logrus"
)

const (
	// Create indicates that a container has been created.
	Create = "create"
	// Destroy indicates that a container has been destroyed.
	Destroy = "destroy"
	// Start indicates that a container has been started.
	Start = "start"
	// Die indicates that a container has been stopped.
	Die = "die"

	containerStoppedMessage = "The container serving this application is not currently running."
)

var evtOptions = types.EventsOptions{
	Filters: filters.NewArgs(),
}

func init() {
	evtOptions.Filters.Add("event", Create)
	evtOptions.Filters.Add("event", Destroy)
	evtOptions.Filters.Add("event", Start)
	evtOptions.Filters.Add("event", Die)
}

// Event represents a particular action that has been taken on a container.
type Event struct {
	Action    string
	Container *container.Container
}

// Dockmon manages a connection to the Docker daemon and delivers events as containers are created, destroyed, started, and stopped. A list of containers is also kept.
type Dockmon struct {
	EventChan  <-chan *Event
	log        *logrus.Entry
	client     *client.Client
	conMap     util.StringMap
	logger     *logger.Logger
	eventChan  chan *Event
	closeFunc  context.CancelFunc
	closedChan chan bool
}

func (d *Dockmon) loggerCallback(resp *http.Response) {
	d.logger.LogResponse(resp)
}

func (d *Dockmon) newContainerFromClient(ctx context.Context, client *client.Client, id string) (*container.Container, error) {
	containerJSON, err := client.ContainerInspect(ctx, id)
	if err != nil {
		return nil, err
	}
	c, err := container.New(
		containerJSON.ID,
		containerJSON.Name[1:],
		containerJSON.Config.Labels,
	)
	if err != nil {
		return nil, err
	}
	if !containerJSON.State.Running {
		c.Disable(containerStoppedMessage)
	}
	c.Proxy.SetCallback(d.loggerCallback)
	return c, nil
}

func (d *Dockmon) sendEvent(action string, c *container.Container) {
	d.eventChan <- &Event{
		Action:    action,
		Container: c,
	}
}

func (d *Dockmon) add(ctx context.Context, id string) {
	if con, err := d.newContainerFromClient(ctx, d.client, id); err == nil {
		d.conMap.Insert(id, con)
		d.sendEvent(Create, con)
		if !con.Disabled {
			d.sendEvent(Start, con)
		}
	}
}

func (d *Dockmon) toggleState(id string, disable bool) {
	if v, ok := d.conMap[id]; ok {
		c := v.(*container.Container)
		if disable {
			c.Disable(containerStoppedMessage)
			d.sendEvent(Die, c)
		} else {
			c.Enable()
			d.sendEvent(Start, c)
		}
	}
}

func (d *Dockmon) sync(ctx context.Context) error {

	// Fetch the list of containers
	d.log.Info("performing sync...")
	containers, err := d.client.ContainerList(
		ctx,
		types.ContainerListOptions{All: true},
	)
	if err != nil {
		return err
	}

	// Create a StringMap of the containers from the new list
	newConMap := util.StringMap{}
	for _, c := range containers {
		newConMap.Insert(c.ID, &container.Container{
			Disabled: c.State != "running",
		})
	}

	// Compare the new map with the old one
	// - toggle the state of any containers that have changed
	// - add any containers that we didn't know were started
	for conID, v := range newConMap {
		oldV, ok := d.conMap[conID]
		if ok {
			c := oldV.(*container.Container)
			if v.(*container.Container).Disabled != c.Disabled {
				d.toggleState(conID, c.Disabled)
			}
		} else {
			d.add(ctx, conID)
		}
	}

	// Remove the containers that we don't know were stopped
	for conID, v := range d.conMap.Difference(newConMap) {
		d.conMap.Remove(conID)
		d.sendEvent(Destroy, v.(*container.Container))
	}

	return nil
}

func (d *Dockmon) loop(ctx context.Context) error {
	msgChan, errChan := d.client.Events(ctx, evtOptions)
	if err := d.sync(ctx); err != nil {
		return err
	}
	for {
		select {
		case msg := <-msgChan:
			switch msg.Action {
			case Create:
				d.add(ctx, msg.ID)
			case Destroy:
				if v, ok := d.conMap[msg.ID]; ok {
					d.conMap.Remove(msg.ID)
					d.sendEvent(Destroy, v.(*container.Container))
				}
			case Start:
				d.toggleState(msg.ID, false)
			case Die:
				d.toggleState(msg.ID, true)
			}
		case err := <-errChan:
			return err
		}
	}
}

func (d *Dockmon) run(ctx context.Context) {
	defer close(d.closedChan)
	defer close(d.eventChan)
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
		eventChan       = make(chan *Event)
		d               = &Dockmon{
			EventChan:  eventChan,
			log:        logrus.WithField("context", "dockmon"),
			client:     c,
			conMap:     util.StringMap{},
			logger:     cfg.Logger,
			eventChan:  eventChan,
			closeFunc:  cancelFunc,
			closedChan: make(chan bool),
		}
	)
	go d.run(ctx)
	return d, nil
}

// StartContainer attempts to start the specified container.
func (d *Dockmon) StartContainer(ctx context.Context, id string) error {
	return d.client.ContainerStart(
		ctx,
		id,
		types.ContainerStartOptions{},
	)
}

// StopContainer attempts to stop the specified container.
func (d *Dockmon) StopContainer(ctx context.Context, id string) error {
	return d.client.ContainerStop(
		ctx,
		id,
		nil,
	)
}

// Close shuts down the connection to the Docker daemon.
func (d *Dockmon) Close() {
	d.closeFunc()
	<-d.closedChan
}
