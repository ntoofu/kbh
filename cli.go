package main

import (
	"io"

	"github.com/urfave/cli"
)

type subcommandsFactory interface {
	GenerateCreateTask() func(*cli.Context) error
	GenerateUpdateTask() func(*cli.Context) error
	GenerateShowTask() func(*cli.Context) error
}

type normalSubcommandsFactory struct {
	stdout io.Writer
}

func generateCliApp(stdout io.Writer, stderr io.Writer, exiter func(int), factory subcommandsFactory) *cli.App {
	cli.OsExiter = exiter
	cli.ErrWriter = stderr
	app := cli.NewApp()
	app.Name = "Kanban Boards Handler"
	app.Usage = "handle many kanban boards or issue trackers in a persistent way"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "config file",
			Value: "config.yml",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "create a new task",
			Action:  factory.GenerateCreateTask(),
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "update state of a task",
			Action:  factory.GenerateUpdateTask(),
		},
		{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "display a list of tasks",
			Action:  factory.GenerateShowTask(),
		},
	}
	app.Writer = stdout
	app.ErrWriter = stderr
	return app
}
