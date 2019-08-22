package main

import (
	"log"
	"net/http"

	"github.com/namsral/flag"
	"gitlab.com/dj80hd/observ/pkg/app"
)

var addr string
var workers int

func init() {
	flag.StringVar(&addr,
		"addr",
		":8111",
		"Address to listen on e.g. :8111",
	)
	flag.IntVar(&workers,
		"workers",
		4,
		"Number of workers e.g. 11",
	)
}

func main() {
	flag.Parse()
	log.Printf("starting on addr %s with %d workers", addr, workers)
	log.Fatal(http.ListenAndServe(addr, app.New(workers)))
}
