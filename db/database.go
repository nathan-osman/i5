package db

// Database provides a generic interface for working with a database server.
type Database interface {

	// Name returns the unique identifier for the database server.
	Name() string

	// Version returns the version string for the database server.
	Version() (string, error)

	// CreateUser ensures that the specified user exists, creating them if not using the provided password.
	CreateUser(username, password string) error

	// CreateDatabase ensures that the specified database exists, creating it if not using the provided owner.
	CreateDatabase(name, user string) error
}
