package main

import (
	// "encoding/json"
	// "fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/medhansh-32/api-gateway/internal/config"
	"github.com/medhansh-32/api-gateway/internal/database"
	"github.com/medhansh-32/api-gateway/internal/handlers"
	"github.com/medhansh-32/api-gateway/internal/repository"
	"github.com/medhansh-32/api-gateway/internal/service"
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

	// b, _ := json.MarshalIndent(gatewayCfg, "", "  ")
	// fmt.Println(string(b))

	port := cfg.ServerPort
	address := "localhost:" + strconv.Itoa(port)
	router := http.DefaultServeMux
	server := http.Server{
		Addr:    address,
		Handler: router,
	}

	

	serveChan := make(chan struct{})
	
	db,dbError := database.NewMysqlConnection(cfg)

	if dbError!=nil{
		log.Fatal("Error Making Connection with Database : ",dbError.Error())
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	// user,err:=userService.GetUserById(1)
	jwtService:= service.NewJWTService(cfg.JWTSecret)
	authService:= service.NewAuthService(userService,jwtService)
	proxyService:= service.NewProxyService()
	gateWayHandler := handlers.NewGateWayHandler(proxyService,*authService)
	gateWayHandler.RegisteRoutes(router)
	
	if err!=nil{
		println(err.Error())
	}
	
	// u, _ := json.MarshalIndent(user, "", "  ")
	// fmt.Println(string(u))
	

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			log.Fatal("Error Starting Server : ", err.Error())
		}
	}()
	log.Print("API GateWay Started Port :", 8080)

	<-serveChan


}
