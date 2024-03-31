package main

import (
	"flotta-home/mindbond/websocket-server/pkg/client"
	"flotta-home/mindbond/websocket-server/pkg/config"
	"flotta-home/mindbond/websocket-server/pkg/server"
	"log"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	authClient := client.InitAuthServiceClient(config.AuthServiceUrl)
	wsServer := server.Server{
		Port:       config.Port,
		AuthClient: authClient,
	}

	wsServer.Start()
}
