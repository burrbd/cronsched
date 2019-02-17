package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/burrbd/cronsched"
)

func main() {
	timeArg := &cronsched.TimeVal{}

	flag.Var(timeArg, "time", "set the time in HH:MM format")
	flag.Parse()

	if !timeArg.IsSet {
		timeArg.Time = time.Now()
	}

	if err := cronsched.Run(os.Stdin, os.Stdout, timeArg.Time); err != nil {
		log.Fatal(err)
	}
}
