package db

// Manager manages connections to database servers.
type Manager struct {
	databases map[string]Database
}

// NewManager creates a new database manager.
func NewManager() *Manager {
	return &Manager{
		databases: map[string]Database{},
	}
}

// Register adds the database to the manager.
func (m *Manager) Register(database Database) {
	m.databases[database.Name()] = database
}
