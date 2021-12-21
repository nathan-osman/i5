package dbman

import (
	"fmt"
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

// Get attempts to retrieve the specified database driver.
func (m *Manager) Get(name string) (Database, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if d, ok := m.databases[name]; ok {
		return d, nil
	}
	return nil, fmt.Errorf("invalid database driver \"%s\" specifed", name)
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
