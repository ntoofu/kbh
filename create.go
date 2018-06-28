package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func (f normalSubcommandsFactory) GenerateCreateTask() func(*cli.Context) error {
	return createTask
}

func createTask(c *cli.Context) error {
	if len(c.Args()) > 3 {
		return fmt.Errorf("too many arguments are given: %v", c.Args()[3:])
	} else if len(c.Args()) < 3 {
		return fmt.Errorf("too few arguments are given: %v", c.Args())
	}

	// TODO
	return fmt.Errorf("Sorry, not implemented yet ...")
}
