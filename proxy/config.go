package proxy

// Config provides the configuration for a proxy instance.
type Config struct {
	// Addr provides the remote address to use for proxying requests to.
	Addr string
	// Mountpoints provides a list of static mountpoints for this proxy.
	Mountpoints []*Mountpoint
}
