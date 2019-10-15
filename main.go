package main

import (
	"github.com/Brialius/calendar/cmd"
	"log"
)

var (
	version = "dev"
	build   = "local"
)

func main() {
	log.Printf("Started calendar %s-%s", version, build)

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
