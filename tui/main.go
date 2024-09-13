package main

import (
	"log"

	"github.com/arnavsurve/taskman/tui/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
