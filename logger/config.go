package logger

// Config provides the configuration for the logger.
type Config struct {
	// Debug indicates whether debugging is enabled or not.
	Debug bool
	// GeolocationDBType indicates the type of geolocation database.
	GeolocationDBType string
	// GeolocationDBPath provides the path to the geolocation database.
	GeolocationDBPath string
	// Status indicates whether the status website is enabled or not.
	Status bool
}
