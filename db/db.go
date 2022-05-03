package db

import (
	"path"

	bolt "go.etcd.io/bbolt"
)

const databaseName = "i5.db"

// DB provides a centralized entrypoint to the Bolt database.
type DB struct {
	db *bolt.DB
}

// New creates and initializes the database.
func New(cfg *Config) (*DB, error) {
	d, err := bolt.Open(
		path.Join(cfg.StorageDir, databaseName),
		0600,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &DB{d}, nil
}

// Close frees all database resources.
func (d *DB) Close() {
	d.db.Close()
}
