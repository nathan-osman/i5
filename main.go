package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/howeyc/gopass"
	"github.com/nathan-osman/i5/conman"
	"github.com/nathan-osman/i5/dbman"
	"github.com/nathan-osman/i5/dockmon"
	"github.com/nathan-osman/i5/geolocation"
	"github.com/nathan-osman/i5/notifier"
	"github.com/nathan-osman/i5/server"
	"github.com/nathan-osman/i5/status"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "i5",
		Usage: "reverse proxy for Docker containers",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				EnvVars: []string{"DEBUG"},
				Usage:   "enable debug mode",
			},
			&cli.StringFlag{
				Name:    "docker-host",
				Value:   "unix:///var/run/docker.sock",
				EnvVars: []string{"DOCKER_HOST"},
				Usage:   "host running Docker",
			},
			&cli.StringFlag{
				Name:    "email",
				EnvVars: []string{"EMAIL"},
				Usage:   "email address to use for challenges",
			},
			&cli.StringFlag{
				Name:    "geolocation-db-type",
				EnvVars: []string{"GEOLOCATION_DB_TYPE"},
				Usage:   "type of IP gelocation database",
			},
			&cli.StringFlag{
				Name:    "geolocation-db-path",
				EnvVars: []string{"GEOLOCATION_DB_PATH"},
				Usage:   "path to IP geolocation database",
			},
			&cli.StringFlag{
				Name:    "http-addr",
				Value:   ":http",
				EnvVars: []string{"HTTP_ADDR"},
				Usage:   "HTTP address to listen on",
			},
			&cli.StringFlag{
				Name:    "https-addr",
				Value:   ":https",
				EnvVars: []string{"HTTPS_ADDR"},
				Usage:   "HTTPS address to listen on",
			},
			&cli.BoolFlag{
				Name:    "mysql",
				EnvVars: []string{"MYSQL"},
				Usage:   "connect to MySQL",
			},
			&cli.IntFlag{
				Name:    "mysql-port",
				Value:   3306,
				EnvVars: []string{"MYSQL_PORT"},
				Usage:   "port for MySQL server",
			},
			&cli.StringFlag{
				Name:    "mysql-host",
				Value:   "mysql",
				EnvVars: []string{"MYSQL_HOST"},
				Usage:   "hostname of MySQL server",
			},
			&cli.StringFlag{
				Name:    "mysql-user",
				Value:   "root",
				EnvVars: []string{"MYSQL_USER"},
				Usage:   "username for connecting to MySQL",
			},
			&cli.StringFlag{
				Name:    "mysql-password",
				EnvVars: []string{"MYSQL_PASSWORD"},
				Usage:   "password for connecting to MySQL",
			},
			&cli.BoolFlag{
				Name:    "postgres",
				EnvVars: []string{"POSTGRES"},
				Usage:   "connect to PostgreSQL",
			},
			&cli.IntFlag{
				Name:    "postgres-port",
				Value:   5432,
				EnvVars: []string{"POSTGRES_PORT"},
				Usage:   "port for PostgreSQL server",
			},
			&cli.StringFlag{
				Name:    "postgres-host",
				Value:   "postgres",
				EnvVars: []string{"POSTGRES_HOST"},
				Usage:   "hostname of PostgreSQL server",
			},
			&cli.StringFlag{
				Name:    "postgres-user",
				Value:   "postgres",
				EnvVars: []string{"POSTGRES_USER"},
				Usage:   "username for connecting to PostgreSQL",
			},
			&cli.StringFlag{
				Name:    "postgres-password",
				EnvVars: []string{"POSTGRES_PASSWORD"},
				Usage:   "password for connecting to PostgreSQL",
			},
			&cli.StringFlag{
				Name:    "status-domain",
				EnvVars: []string{"STATUS_DOMAIN"},
				Usage:   "domain name for internal server",
			},
			&cli.StringFlag{
				Name:    "status-key",
				EnvVars: []string{"STATUS_KEY"},
				Usage:   "secret key for encoding cookies",
			},
			&cli.BoolFlag{
				Name:    "status-insecure",
				EnvVars: []string{"STATUS_INSECURE"},
				Usage:   "allow insecure connections to the internal server",
			},
			&cli.StringFlag{
				Name:    "storage-dir",
				EnvVars: []string{"STORAGE_DIR"},
				Usage:   "directory for storing files related to i5",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "createuser",
				Usage: "create a new user account",
				Action: func(c *cli.Context) error {

					// Prompt for the username
					var username string
					fmt.Print("Username? ")
					fmt.Scanln(&username)

					// Prompt for the password, hiding the input
					fmt.Print("Password? ")
					p, err := gopass.GetPasswd()
					if err != nil {
						return err
					}

					// Create the user
					if err := status.CreateUser(
						c.String("storage-dir"),
						username,
						string(p),
					); err != nil {
						return nil
					}

					fmt.Println("New user created successfully.")

					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {

			// If an IP geolocation database was provided, load it
			var (
				geoDBType = c.String("geolocation-db-type")
				geoDB     *geolocation.Geolocation
			)
			if geoDBType != "" {
				g, err := geolocation.New(&geolocation.Config{
					DBType: geoDBType,
					DBPath: c.String("geolocation-db-path"),
				})
				if err != nil {
					return err
				}
				geoDB = g
				defer g.Close()
			}

			// Check if the status website was enabled
			var (
				statusDomain = c.String("status-domain")
				n            *notifier.Notifier
			)
			if statusDomain != "" {
				n = notifier.New(&notifier.Config{
					Debug: c.Bool("debug"),
				})
				defer n.Close()
			}

			// Create the Docker monitor
			dm, err := dockmon.New(&dockmon.Config{
				Host:        c.String("docker-host"),
				Geolocation: geoDB,
				Notifier:    n,
			})
			if err != nil {
				return err
			}
			defer dm.Close()

			// Create the database manager
			d := dbman.NewManager()

			// Connect to MySQL if requested
			if c.Bool("mysql") {
				msql, err := dbman.NewMySQL(
					c.String("mysql-host"),
					c.Int("mysql-port"),
					c.String("mysql-user"),
					c.String("mysql-password"),
				)
				if err != nil {
					return err
				}
				d.Register(msql)
			}

			// Connect to PostgreSQL if requested
			if c.Bool("postgres") {
				psql, err := dbman.NewPostgres(
					c.String("postgres-host"),
					c.Int("postgres-port"),
					c.String("postgres-user"),
					c.String("postgres-password"),
				)
				if err != nil {
					return err
				}
				d.Register(psql)
			}

			// Create the container manager
			cm := conman.New(&conman.Config{
				EventChan: dm.EventChan,
				Dbman:     d,
			})
			defer cm.Close()

			// If a domain name for the internal server was specified, use it
			if statusDomain != "" {
				s, err := status.New(&status.Config{
					Key:        c.String("status-key"),
					Debug:      c.Bool("debug"),
					Domain:     statusDomain,
					Insecure:   c.Bool("status-insecure"),
					StorageDir: c.String("storage-dir"),
					Conman:     cm,
					Dbman:      d,
					Notifier:   n,
				})
				if err != nil {
					return err
				}
				defer s.Close()
				cm.Add(s.Container)
			}

			// Create the server
			sv, err := server.New(&server.Config{
				Debug:      c.Bool("debug"),
				Email:      c.String("email"),
				HTTPAddr:   c.String("http-addr"),
				HTTPSAddr:  c.String("https-addr"),
				StorageDir: c.String("storage-dir"),
				Conman:     cm,
				Notifier:   n,
			})
			if err != nil {
				return err
			}
			defer sv.Close()

			// Wait for SIGINT or SIGTERM
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			<-sigChan

			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err.Error())
	}
}
