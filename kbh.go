package main

import (
	"os"
)

func main() {
	app := generateCliApp(os.Stdout, os.Stderr, os.Exit, normalSubcommandsFactory{os.Stdout})
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
