package service

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

// Service represents configuration for an application running within a Docker container. The configuration is generated from the container's labels.
type Service struct {
	// ID is the service's unique identifier.
	ID string
	// Name is the service's name.
	Name string
	// Addr is the address of the service's HTTP server.
	Addr string
	// Domains provides a list of domain names for the service.
	Domains []string
	// Insecure indicates that non-TLS traffic should not be upgraded.
	Insecure bool
}

func (s *Service) parseLabels(labels map[string]string) error {
	if addr, ok := labels[labelAddr]; ok {
		s.Addr = addr
	} else {
		return errMissingAddress
	}
	if domains, ok := labels[labelDomains]; ok {
		for _, domain := range strings.Split(domains, ",") {
			s.Domains = append(s.Domains, strings.TrimSpace(domain))
		}
	} else {
		return errMissingDomains
	}
	if insecureStr, ok := labels[labelInsecure]; ok {
		if insecure, err := strconv.ParseBool(insecureStr); err == nil {
			s.Insecure = insecure
		} else {
			return err
		}
	}
	return nil
}

// New creates a new Service from the provided data.
func New(id, name string, labels map[string]string) (*Service, error) {
	s := &Service{
		ID:   id,
		Name: name,
	}
	if err := s.parseLabels(labels); err != nil {
		return nil, err
	}
	return s, nil
}

// NewFromContainer creates a new Service from the specified container.
func NewFromContainer(ctx context.Context, client *client.Client, id string) (*Service, error) {
	if containerJSON, err := client.ContainerInspect(ctx, id); err == nil {
		return New(
			containerJSON.ID,
			containerJSON.Name,
			containerJSON.Config.Labels,
		)
	} else {
		return nil, err
	}
}
