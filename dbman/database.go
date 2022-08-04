package dbman

// Database provides a generic interface for working with a database server.
type Database interface {

	// Name returns the unique identifier for the database server.
	Name() string

	// Title returns the human-readable name of the database server.
	Title() string

	// Version returns the version string for the database server.
	Version() string

	// CreateUser ensures that the specified user exists, creating them if not using the provided password.
	CreateUser(username, password string) error

	// CreateDatabase ensures that the specified database exists, creating it if not using the provided owner.
	CreateDatabase(name, user string) error

	// ListDatabases returns a list of databases in the database server.
	ListDatabases() ([]string, error)

	// DeleteDatabase removes the specified database.
	DeleteDatabase(name string) error

	// Shut down the connection to the database server.
	Close()
}
