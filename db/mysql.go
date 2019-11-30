package db

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

const NameMySQL = "mysql"

// MySQL provides access to a MySQL database.
type MySQL struct {
	conn    *sql.DB
	log     *logrus.Entry
	version string
}

// NewMySQL attempts to create a connection to a MySQL database.
func NewMySQL(host string, port int, user, password string) (*MySQL, error) {
	cfg := mysql.NewConfig()
	cfg.User = user
	cfg.Passwd = password
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%s:%d", host, port)
	cfg.InterpolateParams = true
	c, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	m := &MySQL{
		conn: c,
		log:  logrus.WithField("context", "mysql"),
	}
	if err := m.conn.QueryRow("SELECT VERSION()").
		Scan(&m.version); err != nil {
		return nil, err
	}
	m.log.Info("connected to MySQL")
	return m, nil
}

func (m *MySQL) Name() string {
	return NameMySQL
}

func (m *MySQL) Title() string {
	return "MySQL"
}

func (m *MySQL) Version() string {
	return m.version
}

func (m *MySQL) CreateUser(user, password string) error {
	if _, err := m.conn.Query(
		"CREATE USER IF NOT EXISTS ? IDENTIFIED BY ?",
		user,
		password,
	); err != nil {
		return err
	}
	return nil
}

func (m *MySQL) CreateDatabase(name, user string) error {
	if _, err := m.conn.Query(
		fmt.Sprintf(
			"CREATE DATABASE IF NOT EXISTS `%s`",
			name,
		),
	); err != nil {
		return err
	}
	if _, err := m.conn.Query(
		fmt.Sprintf(
			"GRANT ALL PRIVILEGES ON `%s`.* TO ?",
			name,
		),
		user,
	); err != nil {
		return err
	}
	return nil
}

func (m *MySQL) Close() {
	m.conn.Close()
	m.log.Info("disconnected from MySQL")
}
