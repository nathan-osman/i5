package geolocation

// Config provides the configuration for the geolocation service.
type Config struct {
	// DBType indicates the type of geolocation database.
	DBType string
	// DBPath provides the path to a geolocation database.
	DBPath string
}
