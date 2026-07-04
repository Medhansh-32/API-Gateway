package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/medhansh-32/api-gateway/internal/config"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error Occured while loading Config " + err.Error())
	}
	println("Loaded the API Gateway configuration from " + cfg.RoutingConfig)

	gatewayCfg, err := config.LoadGateWayConfig(cfg.RoutingConfig)
	if err != nil {
		log.Fatal("Error Occured while loading Config " + err.Error())
	}
	println("Loaded the API Gateway configuration from " + gatewayCfg.Gateway.Name)

	b, _ := json.MarshalIndent(gatewayCfg, "", "  ")
	fmt.Println(string(b))

	port := cfg.ServerPort
	address := "localhost:" + strconv.Itoa(port)
	router := http.DefaultServeMux
	server := http.Server{
		Addr:    address,
		Handler: router,
	}

	serveChan := make(chan struct{})
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			log.Fatal("Error Starting Server : ", err.Error())
		}
	}()
	log.Print("Server Started Port :", 8080)

	<-serveChan


}
