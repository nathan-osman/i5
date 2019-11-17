package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Postgres provides access to a PostgreSQL database.
type Postgres struct {
	conn *sql.DB
}

// NewPostgres attempts to create a connection to a PostgreSQL database.
func NewPostgres(host string, port int, user, password string) (*Postgres, error) {
	c, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s sslmode=disable",
			host,
			port,
			user,
			password,
		),
	)
	if err != nil {
		return nil, err
	}
	return &Postgres{
		conn: c,
	}, nil
}

func (p *Postgres) CreateUser(user, password string) error {
	r, err := p.conn.Query("SELECT 1 FROM pg_roles WHERE rolname=$1", user)
	if err != nil {
		return err
	}
	defer r.Close()
	if !r.Next() {
		if _, err := p.conn.Query("CREATE USER $1 WITH PASSWORD $2"); err != nil {
			return err
		}
	}
	return nil
}

func (p *Postgres) CreateDatabase(name, user string) error {
	r, err := p.conn.Query("SELECT 1 FROM pg_database WHERE datname=$1", name)
	if err != nil {
		return err
	}
	defer r.Close()
	if !r.Next() {
		if _, err := p.conn.Query("CREATE DATABASE $1 WITH OWNER $2", name, user); err != nil {
			return err
		}
	}
	return nil
}
