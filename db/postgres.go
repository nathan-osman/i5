package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const NamePostgreSQL = "postgres"

// Postgres provides access to a PostgreSQL database.
type Postgres struct {
	conn    *sql.DB
	log     *logrus.Entry
	version string
}

// NewPostgres attempts to create a connection to a PostgreSQL database. The server version is retrieved as well.
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
	p := &Postgres{
		conn: c,
		log:  logrus.WithField("context", "conman"),
	}
	if err := p.conn.QueryRow("SHOW server_version").
		Scan(&p.version); err != nil {
		return nil, err
	}
	p.log.Info("connected to PostgreSQL")
	return p, nil
}

func (p *Postgres) Name() string {
	return NamePostgreSQL
}

func (p *Postgres) Title() string {
	return "PostgreSQL"
}

func (p *Postgres) Version() string {
	return p.version
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

func (p *Postgres) Close() {
	p.conn.Close()
	p.log.Info("disconnected from PostgreSQL")
}
