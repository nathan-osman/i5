package dockmon

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/nathan-osman/i5/geolocation"
	"github.com/nathan-osman/i5/notifier"
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
	Container *Container
}

// Dockmon manages a connection to the Docker daemon and delivers events as containers are created, destroyed, started, and stopped. A list of containers is also kept.
type Dockmon struct {
	EventChan   <-chan *Event
	log         *logrus.Entry
	client      *client.Client
	conMap      util.StringMap
	geolocation *geolocation.Geolocation
	notifier    *notifier.Notifier
	eventChan   chan *Event
	closeFunc   context.CancelFunc
	closedChan  chan bool
}

func (d *Dockmon) sendEvent(action string, con *Container) {
	d.eventChan <- &Event{
		Action:    action,
		Container: con,
	}
}

func (d *Dockmon) add(ctx context.Context, id string) {
	if con, err := d.newContainerFromClient(ctx, d.client, id); err == nil {
		d.conMap.Insert(id, con)
		d.sendEvent(Create, con)
		if con.Running {
			d.sendEvent(Start, con)
		}
	}
}

func (d *Dockmon) toggleState(id string, running bool) {
	if v, ok := d.conMap[id]; ok {
		con := v.(*Container)
		con.Running = running
		if running {
			d.sendEvent(Start, con)
		} else {
			d.sendEvent(Die, con)
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
	for _, con := range containers {
		newConMap.Insert(con.ID, &Container{
			Running: con.State == "running",
		})
	}

	// Compare the new map with the old one
	// - toggle the state of any containers that have changed
	// - add any containers that we didn't know were started
	for conID, v := range newConMap {
		oldV, ok := d.conMap[conID]
		if ok {
			con := oldV.(*Container)
			if v.(*Container).Running != con.Running {
				d.toggleState(conID, con.Running)
			}
		} else {
			d.add(ctx, conID)
		}
	}

	// Remove the containers that we don't know were stopped
	for conID, v := range d.conMap.Difference(newConMap) {
		d.conMap.Remove(conID)
		d.sendEvent(Destroy, v.(*Container))
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
					d.sendEvent(Destroy, v.(*Container))
				}
			case Start:
				d.toggleState(msg.ID, true)
			case Die:
				d.toggleState(msg.ID, false)
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
			EventChan:   eventChan,
			log:         logrus.WithField("context", "dockmon"),
			client:      c,
			conMap:      util.StringMap{},
			geolocation: cfg.Geolocation,
			notifier:    cfg.Notifier,
			eventChan:   eventChan,
			closeFunc:   cancelFunc,
			closedChan:  make(chan bool),
		}
	)
	go d.run(ctx)
	return d, nil
}

// Close shuts down the connection to the Docker daemon.
func (d *Dockmon) Close() {
	d.closeFunc()
	<-d.closedChan
}
