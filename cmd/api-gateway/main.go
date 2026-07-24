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
	"github.com/medhansh-32/api-gateway/internal/middleware"
	"github.com/medhansh-32/api-gateway/internal/repository"
	"github.com/medhansh-32/api-gateway/internal/schedular"
	"github.com/medhansh-32/api-gateway/internal/service"
)

func main() {

	cfg, err := config.LoadApplicationConfig()
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
	

	serveChan := make(chan struct{})
	
	db,dbError := database.NewMysqlConnection(cfg)

	if dbError!=nil{
		log.Fatal("Error Making Connection with Database : ",dbError.Error())
	}

	cfgManager := &config.ConfigManager{}
	
	cfgManager.UpdateGateWayConfig(gatewayCfg)

	
	configWatcher := config.NewConfigWatcher(cfgManager)

	go configWatcher.WatchGateWayConfig(cfg.RoutingConfig)
	go configWatcher.WatchApplicationConfig(".env")

	healthService := service.NewHealthService(cfgManager)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	// user,err:=userService.GetUserById(1)
	jwtService:= service.NewJWTService(cfg.JWTSecret)
	authService:= service.NewAuthService(userService,jwtService)
	proxyService:= service.NewProxyService(cfgManager,healthService)
	gateWayHandler := handlers.NewGateWayHandler(proxyService,*authService)


	auth := middleware.NewAuthMiddleWare(authService,cfgManager)
	rateLimit := middleware.NewRateLimitingMiddleware(cfgManager)
	
	router := http.DefaultServeMux
	middlewareRouter := middleware.Logger(auth.Authentication(rateLimit.RateLimitCheck(router)))

	port := cfg.ServerPort
	address := "localhost:" + strconv.Itoa(port)

	server := http.Server{
		Addr:    address,
		Handler: middlewareRouter,
	}

	gateWayHandler.RegisteRoutes(router)
	
	if err!=nil{
		println(err.Error())
	}
	
	healthCheckSchedular := schedular.NewHealthCheckSchedular(healthService)
	go healthCheckSchedular.InitHealthCheckSchedular()

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
