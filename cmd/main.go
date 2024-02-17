package main

import (
	"flotta-home/mindbond/websocket-server/pkg/config"
	"flotta-home/mindbond/websocket-server/pkg/server"
	"log"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}
	server.Start(config.Port)
}
