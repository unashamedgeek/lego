package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

// Logger is used to log errors; if nil, the default log.Logger is used.
var Logger *log.Logger

// logger is an helper function to retrieve the available logger
func logger() *log.Logger {
	if Logger == nil {
		Logger = log.New(os.Stderr, "", log.LstdFlags)
	}
	return Logger
}

func main() {

	app := cli.NewApp()
	app.Name = "lego"
	app.Usage = "Let's encrypt client to go!"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Create and install a certificate",
			Action: run,
		},
		{
			Name:  "auth",
			Usage: "Create a certificate",
			Action: func(c *cli.Context) {
				logger().Fatal("Not implemented")
			},
		},
		{
			Name:  "install",
			Usage: "Install a certificate",
			Action: func(c *cli.Context) {
				logger().Fatal("Not implemented")
			},
		},
		{
			Name:  "revoke",
			Usage: "Revoke a certificate",
			Action: func(c *cli.Context) {
				logger().Fatal("Not implemented")
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "certificate",
					Usage: "Revoke a specific certificate",
				},
				cli.StringFlag{
					Name:  "key",
					Usage: "Revoke all certs generated by the provided authorized key.",
				},
			},
		},
		{
			Name:  "rollback",
			Usage: "Rollback a certificate",
			Action: func(c *cli.Context) {
				logger().Fatal("Not implemented")
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "checkpoints",
					Usage: "Revert configuration N number of checkpoints",
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "domains, d",
			Usage: "Add domains to the process",
		},
		cli.StringFlag{
			Name:  "server, s",
			Value: "https://www.letsencrypt-demo.org/acme/new-reg",
			Usage: "CA hostname (and optionally :port). The server certificate must be trusted in order to avoid further modifications to the client.",
		},
		cli.StringFlag{
			Name:  "authkey, k",
			Usage: "Path to the authorized key file",
		},
		cli.StringFlag{
			Name:  "email, m",
			Usage: "Email used for registration and recovery contact.",
		},
		cli.IntFlag{
			Name:  "rsa-key-size, B",
			Value: 2048,
			Usage: "Size of the RSA key.",
		},
		cli.BoolFlag{
			Name:  "no-confirm",
			Usage: "Turn off confirmation screens.",
		},
		cli.BoolFlag{
			Name:  "agree-tos, e",
			Usage: "Skip the end user license agreement screen.",
		},
		cli.StringFlag{
			Name:  "config-dir",
			Value: configDir,
			Usage: "Configuration directory.",
		},
		cli.StringFlag{
			Name:  "work-dir",
			Value: workDir,
			Usage: "Working directory.",
		},
		cli.StringFlag{
			Name:  "backup-dir",
			Value: backupDir,
			Usage: "Configuration backups directory.",
		},
		cli.StringFlag{
			Name:  "key-dir",
			Value: keyDir,
			Usage: "Keys storage.",
		},
		cli.StringFlag{
			Name:  "cert-dir",
			Value: certDir,
			Usage: "Certificates storage.",
		},
	}

	app.Run(os.Args)
}