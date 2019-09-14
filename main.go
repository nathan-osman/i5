//go:generate go run assets/generate.go

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nathan-osman/i5/conman"
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
		cli.StringFlag{
			Name:   "status-domain",
			EnvVar: "STATUS_DOMAIN",
			Usage:  "domain name for internal server",
		},
		cli.StringFlag{
			Name:   "storage-dir",
			EnvVar: "STORAGE_DIR",
			Usage:  "directory for storing certificates",
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

		// Create the container manager
		cm := conman.New(&conman.Config{
			ConStartedChan: dm.ConStartedChan,
			ConStoppedChan: dm.ConStoppedChan,
		})
		defer cm.Close()

		// If a domain name for the internal server was specified, use it
		if statusDomain := c.String("status-domain"); statusDomain != "" {
			cm.Add(status.New(&status.Config{
				Domain: statusDomain,
				Conman: cm,
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
