package main

import (
	"fmt"
	"os"
)

func main() {
	app := generateCliApp(os.Stdout, os.Stderr, os.Exit, normalSubcommandsFactory{os.Stdout})
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
