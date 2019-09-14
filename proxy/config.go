package proxy

// Config provides the configuration for a proxy instance.
type Config struct {
	Addr        string
	Mountpoints []*Mountpoint
}
