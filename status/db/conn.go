package db

import (
	"path"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Conn maintains a connection to the database.
type Conn struct {
	*gorm.DB
}

// New attempts to connect to the database.
func New(storageDir string) (*Conn, error) {
	d, err := gorm.Open(
		sqlite.Open(path.Join(storageDir, "db.sqlite3")),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}
	if err := d.AutoMigrate(
		&User{},
	); err != nil {
		return nil, err
	}
	return &Conn{
		DB: d,
	}, nil
}

// Close closes the database connection.
func (c *Conn) Close() {
	db, _ := c.DB.DB()
	db.Close()
}
