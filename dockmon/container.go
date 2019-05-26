package dockmon

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/docker/docker/client"
)

const (
	labelAddr     = "i5.addr"
	labelDomains  = "i5.domains"
	labelInsecure = "i5.insecure"
)

var (
	errMissingAddress = errors.New("missing address label")
	errMissingDomains = errors.New("missing domains label")
)

// Container represents configuration for an application running within a Docker container. The configuration is generated from the container's labels.
type Container struct {
	// ID is the container's unique identifier.
	ID string
	// Name is the container's name.
	Name string
	// Addr is the address of the container's HTTP server.
	Addr string
	// Domains provides a list of domain names for the container.
	Domains []string
	// Insecure indicates that non-TLS traffic should not be upgraded.
	Insecure bool
}

func (c *Container) parseLabels(labels map[string]string) error {
	if addr, ok := labels[labelAddr]; ok {
		c.Addr = addr
	} else {
		return errMissingAddress
	}
	if domains, ok := labels[labelDomains]; ok {
		for _, domain := range strings.Split(domains, ",") {
			c.Domains = append(c.Domains, strings.TrimSpace(domain))
		}
	} else {
		return errMissingDomains
	}
	if insecureStr, ok := labels[labelInsecure]; ok {
		if insecure, err := strconv.ParseBool(insecureStr); err == nil {
			c.Insecure = insecure
		} else {
			return err
		}
	}
	return nil
}

// New creates a new Container from the provided data.
func NewContainer(id, name string, labels map[string]string) (*Container, error) {
	c := &Container{
		ID:   id,
		Name: name,
	}
	if err := c.parseLabels(labels); err != nil {
		return nil, err
	}
	return c, nil
}

// NewContainerFromClient creates a new Container using the provided client.
func NewContainerFromClient(ctx context.Context, client *client.Client, id string) (*Container, error) {
	if containerJSON, err := client.ContainerInspect(ctx, id); err == nil {
		return NewContainer(
			containerJSON.ID,
			containerJSON.Name[1:],
			containerJSON.Config.Labels,
		)
	} else {
		return nil, err
	}
}
