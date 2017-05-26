package main

import (
	"flag"
	"log"
)

func main() {
	port := flag.Int("p", 8080, "port to run on")
	flag.Parse()

	NewService().serve(*port)
}

func Log(msg string, args ...interface{}) {
	log.Printf(msg, args...)
}
