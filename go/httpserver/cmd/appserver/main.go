package main

import (
	"flag"
	"httpserver-cncamp/internal/appserver"
	_ "net/http/pprof"

	log "k8s.io/klog/v2"
)

func main() {
	// Set up logging
	log.InitFlags(nil)

	// get config
	var address string
	flag.StringVar(&address, "address", ":8080", "HTTP Server Address")
	flag.Set("v", "4")
	flag.Set("logtostderr", "true")
	flag.Parse()

	srv, _ := appserver.NewServer(address)
	err := srv.Run()
	log.Fatalf("Couldn't run: %s", err)
}
