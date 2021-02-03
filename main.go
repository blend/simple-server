package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/blend/simple-server/server"
)

func main() {
	hostPtr := flag.String("host", "127.0.0.1", "The server host")
	portPtr := flag.Int("port", 80, "The server port")
	configPathPtr := flag.String("config", "config.json", "The server config file in JSON format")

	flag.Parse()

	config, err := server.ReadConfig(*configPathPtr)
	if err != nil {
		log.Fatal(err)
	}
	server, err := server.NewServer(fmt.Sprintf("%s:%d", *hostPtr, *portPtr), config)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.ListenAndServe())
}
