package main

import (
	"context"
	"log"
	"os"
	"text/tabwriter"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
	confire "go.rtnl.ai/confire/usage"
	"go.rtnl.ai/uptime/pkg"
	"go.rtnl.ai/uptime/pkg/config"
	"go.rtnl.ai/uptime/pkg/server"
)

func main() {
	godotenv.Load()

	app := cli.Command{
		Name:    "uptime",
		Usage:   "Service status monitor for Rotational applications and systems.",
		Version: pkg.Version(false),
		Commands: []*cli.Command{
			{
				Name:     "serve",
				Usage:    "Start the uptime server",
				Action:   serve,
				Category: "server",
			},
			{
				Name:     "config",
				Usage:    "Print the uptime configuration guide",
				Action:   usage,
				Category: "server",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "print in list mode instead of table mode",
					},
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func serve(ctx context.Context, c *cli.Command) (err error) {
	var srv *server.Server
	if srv, err = server.New(); err != nil {
		return cli.Exit(err, 1)
	}

	if err = srv.Serve(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func usage(ctx context.Context, c *cli.Command) error {
	tabs := tabwriter.NewWriter(os.Stdout, 1, 0, 4, ' ', 0)
	format := confire.DefaultTableFormat
	if c.Bool("list") {
		format = confire.DefaultListFormat
	}

	var conf config.Config
	if err := confire.Usagef(config.Prefix, &conf, tabs, format); err != nil {
		return cli.Exit(err, 1)
	}
	tabs.Flush()
	return nil
}
