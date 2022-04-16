package dockmon

import (
	"github.com/nathan-osman/i5/logger"
)

// Config provides the configuration for the Docker monitor.
type Config struct {
	// Host provides the address to use for connecting to the Docker daemon.
	Host string
	// Logger is a pointer to a Logger instance
	Logger *logger.Logger
}
