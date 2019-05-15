package main

import (
	"log"

	"github.com/faroyam/pic-resizer/server"

	"github.com/faroyam/pic-resizer/resizer"
)

func main() {
	resizer.StartService()
	log.Fatalf("stopped: %v", server.Serve())
}
