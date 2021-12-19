package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/db"
	"github.com/nathan-osman/i5/dockmon"
	"github.com/nathan-osman/i5/server"
	"github.com/nathan-osman/i5/status"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "i5"
	app.Usage = "reverse proxy for Docker containers"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			EnvVar: "DEBUG",
			Usage:  "enable debug mode",
		},
		cli.StringFlag{
			Name:   "docker-host",
			Value:  "unix:///var/run/docker.sock",
			EnvVar: "DOCKER_HOST",
			Usage:  "host running Docker",
		},
		cli.StringFlag{
			Name:   "email",
			EnvVar: "EMAIL",
			Usage:  "email address to use for challenges",
		},
		cli.StringFlag{
			Name:   "http-addr",
			Value:  ":http",
			EnvVar: "HTTP_ADDR",
			Usage:  "HTTP address to listen on",
		},
		cli.StringFlag{
			Name:   "https-addr",
			Value:  ":https",
			EnvVar: "HTTPS_ADDR",
			Usage:  "HTTPS address to listen on",
		},
		cli.BoolFlag{
			Name:   "mysql",
			EnvVar: "MYSQL",
			Usage:  "connect to MySQL",
		},
		cli.IntFlag{
			Name:   "mysql-port",
			Value:  3306,
			EnvVar: "MYSQL_PORT",
			Usage:  "port for MySQL server",
		},
		cli.StringFlag{
			Name:   "mysql-host",
			Value:  "mysql",
			EnvVar: "MYSQL_HOST",
			Usage:  "hostname of MySQL server",
		},
		cli.StringFlag{
			Name:   "mysql-user",
			Value:  "root",
			EnvVar: "MYSQL_USER",
			Usage:  "username for connecting to MySQL",
		},
		cli.StringFlag{
			Name:   "mysql-password",
			EnvVar: "MYSQL_PASSWORD",
			Usage:  "password for connecting to MySQL",
		},
		cli.BoolFlag{
			Name:   "postgres",
			EnvVar: "POSTGRES",
			Usage:  "connect to PostgreSQL",
		},
		cli.IntFlag{
			Name:   "postgres-port",
			Value:  5432,
			EnvVar: "POSTGRES_PORT",
			Usage:  "port for PostgreSQL server",
		},
		cli.StringFlag{
			Name:   "postgres-host",
			Value:  "postgres",
			EnvVar: "POSTGRES_HOST",
			Usage:  "hostname of PostgreSQL server",
		},
		cli.StringFlag{
			Name:   "postgres-user",
			Value:  "postgres",
			EnvVar: "POSTGRES_USER",
			Usage:  "username for connecting to PostgreSQL",
		},
		cli.StringFlag{
			Name:   "postgres-password",
			EnvVar: "POSTGRES_PASSWORD",
			Usage:  "password for connecting to PostgreSQL",
		},
		cli.StringFlag{
			Name:   "status-domain",
			EnvVar: "STATUS_DOMAIN",
			Usage:  "domain name for internal server",
		},
		cli.BoolFlag{
			Name:   "status-insecure",
			EnvVar: "STATUS_INSECURE",
			Usage:  "allow insecure connections to the internal server",
		},
		cli.StringFlag{
			Name:   "storage-dir",
			EnvVar: "STORAGE_DIR",
			Usage:  "directory for storing files related to i5",
		},
	}
	app.Action = func(c *cli.Context) error {

		// Enable debug logging if requested
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}

		// Create the Docker monitor
		dm, err := dockmon.New(&dockmon.Config{
			Host: c.String("docker-host"),
		})
		if err != nil {
			return err
		}
		defer dm.Close()

		// Create the database manager
		dbman := db.NewManager()

		// Connect to MySQL if requested
		if c.Bool("mysql") {
			msql, err := db.NewMySQL(
				c.String("mysql-host"),
				c.Int("mysql-port"),
				c.String("mysql-user"),
				c.String("mysql-password"),
			)
			if err != nil {
				return err
			}
			dbman.Register(msql)
		}

		// Connect to PostgreSQL if requested
		if c.Bool("postgres") {
			psql, err := db.NewPostgres(
				c.String("postgres-host"),
				c.Int("postgres-port"),
				c.String("postgres-user"),
				c.String("postgres-password"),
			)
			if err != nil {
				return err
			}
			dbman.Register(psql)
		}

		// Create the container manager
		cm := conman.New(&conman.Config{
			EventChan: dm.EventChan,
			Dbman:     dbman,
		})
		defer cm.Close()

		// If a domain name for the internal server was specified, use it
		if statusDomain := c.String("status-domain"); statusDomain != "" {
			cm.Add(status.New(&status.Config{
				Domain:   statusDomain,
				Insecure: c.Bool("status-insecure"),
				Conman:   cm,
				Dbman:    dbman,
			}))
		}

		// Create the server
		sv, err := server.New(&server.Config{
			Debug:      c.Bool("debug"),
			Email:      c.String("email"),
			HTTPAddr:   c.String("http-addr"),
			HTTPSAddr:  c.String("https-addr"),
			StorageDir: c.String("storage-dir"),
			Conman:     cm,
		})
		if err != nil {
			return err
		}
		defer sv.Close()

		// Wait for SIGINT or SIGTERM
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		return nil
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err.Error())
	}
}
