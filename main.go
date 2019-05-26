package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nathan-osman/i5/dockmon"
	"github.com/nathan-osman/i5/server"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "i5"
	app.Usage = "reverse proxy for Docker containers"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "docker-host",
			EnvVar: "DOCKER_HOST",
			Usage:  "host running Docker",
		},
		cli.StringFlag{
			Name:   "server-addr",
			Value:  ":http",
			EnvVar: "SERVER_ADDR",
			Usage:  "address to listen on",
		},
	}
	app.Action = func(c *cli.Context) error {

		// Create the Docker monitor
		dm, err := dockmon.New(&dockmon.Config{
			Host: c.String("docker-host"),
		})
		if err != nil {
			return err
		}
		defer dm.Close()

		// Create the server
		sv, err := server.New(&server.Config{
			Addr:    c.String("server-addr"),
			Dockmon: dm,
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
