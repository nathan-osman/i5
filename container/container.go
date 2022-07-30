package container

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nathan-osman/i5/proxy"
)

const (
	labelAddr             = "i5.addr"
	labelDomains          = "i5.domains"
	labelInsecure         = "i5.insecure"
	labelDatabaseDriver   = "i5.database.driver"
	labelDatabaseName     = "i5.database.name"
	labelDatabaseUser     = "i5.database.user"
	labelDatabasePassword = "i5.database.password"
	labelMountpoints      = "i5.mountpoints"
	labelTitle            = "i5.title"
)

var (
	errMissingDomains           = errors.New("missing domains label")
	errInvalidMountpoint        = errors.New("invalid mountpoint")
	errMissingAddrOrMountpoints = errors.New("missing addr or mountpoints")
)

// Database represents database requirements for a container.
type Database struct {
	Driver   string
	Name     string
	User     string
	Password string
}

// ContainerData provides a base type for container data.
type ContainerData struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Domains     []string `json:"domains"`
	Insecure    bool     `json:"insecure"`
	Running     bool     `json:"running"`
	Maintenance bool     `json:"maintenance"`
	Uptime      int64    `json:"uptime"`
}

// Container represents configuration for a Docker container with i5 metadata.
type Container struct {
	ContainerData
	Database *Database
	Handler  http.Handler
	Proxy    *proxy.Proxy
}

func getWithDefault(m map[string]string, key, def string) string {
	if v, ok := m[key]; ok {
		return v
	}
	return def
}

// New creates a new Container instance from the provided labels.
func New(id, name string, labels map[string]string, running bool) (*Container, error) {
	var (
		cfg = &proxy.Config{}
		c   = &Container{
			ContainerData: ContainerData{
				ID:    id,
				Name:  name,
				Title: name,
			},
		}
	)
	if addr, ok := labels[labelAddr]; ok {
		cfg.Addr = addr
	}
	var (
		databaseDriver, _   = labels[labelDatabaseDriver]
		databasePassword, _ = labels[labelDatabasePassword]
	)
	if databaseDriver != "" && databasePassword != "" {
		c.Database = &Database{
			Driver:   databaseDriver,
			Name:     getWithDefault(labels, labelDatabaseName, name),
			User:     getWithDefault(labels, labelDatabaseUser, name),
			Password: databasePassword,
		}
	}
	if domains, ok := labels[labelDomains]; ok {
		for _, domain := range strings.Split(domains, ",") {
			c.Domains = append(c.Domains, strings.TrimSpace(domain))
		}
	} else {
		return nil, errMissingDomains
	}
	if insecureStr, ok := labels[labelInsecure]; ok {
		if insecure, err := strconv.ParseBool(insecureStr); err == nil {
			c.Insecure = insecure
		} else {
			return nil, err
		}
	}
	if mountpoints, ok := labels[labelMountpoints]; ok {
		for _, mountpoint := range strings.Split(mountpoints, ",") {
			if parts := strings.SplitN(mountpoint, "=", 2); len(parts) == 2 {
				cfg.Mountpoints = append(cfg.Mountpoints, &proxy.Mountpoint{
					Path: parts[0],
					Dir:  parts[1],
				})
			} else {
				return nil, errInvalidMountpoint
			}
		}
	}
	if cfg.Addr == "" && cfg.Mountpoints == nil {
		return nil, errMissingAddrOrMountpoints
	}
	if title, ok := labels[labelTitle]; ok {
		c.Title = title
	}
	c.Proxy = proxy.New(cfg)
	if running {
		c.Started()
	} else {
		c.Stopped()
	}
	return c, nil
}

// Started switches the handler to the proxy for the container.
func (c *Container) Started() {
	c.Running = true
	c.Uptime = time.Now().Unix()
	c.Handler = c.Proxy
}

// Stopped switches the handler to the static "stopped" message.
func (c *Container) Stopped() {
	c.Running = false
	c.Uptime = 0
	c.Handler = &disabledHandler{
		Message: stoppedContainerMessage,
	}
}
