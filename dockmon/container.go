package dockmon

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/docker/docker/client"
	"github.com/nathan-osman/i5/proxy"
)

const (
	labelAddr             = "i5.addr"
	labelDatabaseDriver   = "i5.database.driver"
	labelDatabaseName     = "i5.database.name"
	labelDatabaseUser     = "i5.database.user"
	labelDatabasePassword = "i5.database.password"
	labelDomains          = "i5.domains"
	labelInsecure         = "i5.insecure"
	labelMountpoints      = "i5.mountpoints"
)

var (
	errMissingDomains           = errors.New("missing domains label")
	errInvalidMountpoint        = errors.New("invalid mountpoint")
	errMissingAddrOrMountpoints = errors.New("missing addr or mountpoints")
)

// Database represents database requirements for a container.
type Database struct {
	// Driver indicates which database type is needed.
	Driver string
	// Name indicates the name of the database to connect to.
	Name string
	// User indicates the username for connecting to the database.
	User string
	// Password indicates the password for connecting to the database.
	Password string
}

// Container represents configuration for an application in a Docker container. The configuration is generated from the container's labels.
type Container struct {
	// ID is the container's unique identifier.
	ID string
	// Name is the container's name.
	Name string
	// Domains provides a list of domain names for the container.
	Domains []string
	// Insecure indicates that non-TLS traffic should not be upgraded.
	Insecure bool
	// Handler is used for serving content from the container.
	Handler http.Handler
	// Database provides database requirements for the container.
	Database *Database
	// Running indicates whether the container is running or not
	Running bool
}

func getWithDefault(m map[string]string, key, def string) string {
	if v, ok := m[key]; ok {
		return v
	}
	return def
}

// NewContainer creates a new Container from the provided data.
func NewContainer(id, name string, labels map[string]string, running bool) (*Container, error) {
	var (
		cfg = &proxy.Config{}
		c   = &Container{
			ID:      id,
			Name:    name,
			Running: running,
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
	c.Handler = proxy.New(cfg)
	return c, nil
}

// NewContainerFromClient creates a new Container using the provided client.
func NewContainerFromClient(ctx context.Context, client *client.Client, id string) (*Container, error) {
	containerJSON, err := client.ContainerInspect(ctx, id)
	if err != nil {
		return nil, err
	}
	return NewContainer(
		containerJSON.ID,
		containerJSON.Name[1:],
		containerJSON.Config.Labels,
		containerJSON.State.Running,
	)
}
