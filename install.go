package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/urfave/cli/v2"
)

const systemdUnitFile = `[Unit]
Description=i5 Reverse Proxy
Wants=network-online.target
After=network-online.target

[Service]
ExecStart={{.path}}

[Install]
WantedBy=multi-user.target
`

var installCommand = &cli.Command{
	Name:   "install",
	Usage:  "install the application as a local service",
	Action: install,
}

func install(c *cli.Context) error {

	// Determine the full path to the executable
	p, err := os.Executable()
	if err != nil {
		return err
	}

	// Compile the template
	t, err := template.New("").Parse(systemdUnitFile)
	if err != nil {
		return err
	}

	// Attempt to open the file
	f, err := os.Create("/lib/systemd/system/i5.service")
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the template
	t.Execute(f, map[string]interface{}{
		"path": p,
	})

	fmt.Println("Service installed!")
	fmt.Println("")
	fmt.Println("To modify command line arguments, please edit:")
	fmt.Println("  /lib/systemd/system/i5.service")
	fmt.Println("")
	fmt.Println("To enable the service and start it, run:")
	fmt.Println("  systemctl enable lampctl")
	fmt.Println("  systemctl start lampctl")

	return nil
}
