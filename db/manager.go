package db

import (
	"sync"
)

// Manager manages connections to database servers.
type Manager struct {
	mutex     sync.RWMutex
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
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.databases[database.Name()] = database
}

// List returns a slice of all registered databases.
func (m *Manager) List() []Database {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	l := []Database{}
	for _, d := range m.databases {
		l = append(l, d)
	}
	return l
}
