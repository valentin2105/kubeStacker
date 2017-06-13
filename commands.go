package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/valentin2105/kubeStacker/command"
)

var GlobalFlags = []cli.Flag{}

var Commands = []cli.Command{
	{
		Name:   "add",
		Usage:  "",
		Action: command.CmdAdd,

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name, n",
				Value: "default",
				Usage: "Stack name",
			},
			cli.StringFlag{
				Name:  "type, t",
				Usage: "Stack type (Wordpress, Drupal, PHP...)",
			},

			cli.IntFlag{
				Name:  "size, s",
				Usage: "Stack size (in GB)",
			},
		},
	},
	{
		Name:   "show",
		Usage:  "",
		Action: command.CmdShow,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "delete",
		Usage:  "",
		Action: command.CmdDelete,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.\n ", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
