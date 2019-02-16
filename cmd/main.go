package main

import (
	"bufio"
	"log"
	"os"
	"time"

	"github.com/burrbd/cronlendar"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	if err := cronlendar.Run(reader, os.Stdout, time.Now()); err != nil {
		log.Fatal(err)
	}
}
