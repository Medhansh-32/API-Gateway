package main

import (
	"log"

	"github.com/medhansh-32/api-gateway/internal/config"
)

func main() {
	cfg,err := config.Load();
	if err != nil{
		log.Fatal("Error Occured while loading Config "+err.Error())
	}
	println("Loading the API Gateway configuration from "+cfg.RoutingConfig)
}