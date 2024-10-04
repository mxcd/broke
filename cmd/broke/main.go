package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "broke",
		Usage: "Identity broker",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "debug output",
				EnvVars: []string{"BROKE_VERBOSE"},
			},
			&cli.BoolFlag{
				Name:    "very-verbose",
				Aliases: []string{"vv"},
				Usage:   "trace output",
				EnvVars: []string{"BROKE_VERY_VERBOSE"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{"s"},
				Usage:   "Run identity broker",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "",
						Usage:   "*.broke.yml file to be used",
						EnvVars: []string{"BROKE_CONFIG_FILE"},
					},
				},
			},
			{
				Name:  "plan",
				Usage: "Plan the broker run",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "",
						Usage:   "*.broke.yml file to be used",
						EnvVars: []string{"BROKE_CONFIG_FILE"},
					},
				},
				Action: func(c *cli.Context) error {
					initApplication(c)
					// TODO
					return nil
				},
			},
			{
				Name:  "test",
				Usage: "test the connection to all configured APIs",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "",
						Usage:   "*.broke.yml file to be used",
						EnvVars: []string{"BROKE_CONFIG_FILE"},
					},
				},
				Action: func(c *cli.Context) error {
					initApplication(c)
					// TODO
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func initApplication(c *cli.Context) error {
	util.PrintLogo(c)
	util.SetLogLevel(c)
	util.SetCliContext(c)
	util.GetRootDir()
	err := state.LoadState(c)
	return err
}
