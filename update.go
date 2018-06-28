package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func (f normalSubcommandsFactory) GenerateUpdateTask() func(*cli.Context) error {
	return updateTask
}

func updateTask(c *cli.Context) error {
	if len(c.Args()) > 2 {
		return fmt.Errorf("too many arguments are given: %v", c.Args()[2:])
	} else if len(c.Args()) < 2 {
		return fmt.Errorf("too few arguments are given: %v", c.Args())
	}

	// TODO
	return fmt.Errorf("Sorry, not implemented yet ...")
}
