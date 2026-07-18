package service

import (
	"errors"
	"net/http"
	"strings"

	"github.com/medhansh-32/api-gateway/internal/config"
	"github.com/medhansh-32/api-gateway/internal/models"
)


type ProxyService interface{
	FindTargetRouteForRequest(request *http.Request) (*models.TargetRoute,error)
}

type ProxyServiceImpl struct{
	ConfigManager *config.ConfigManager
}

func NewProxyService(Cfg *config.ConfigManager) (ProxyService){
	return &ProxyServiceImpl{ConfigManager: Cfg}
}

func (proxyServiceImpl *ProxyServiceImpl) FindTargetRouteForRequest(request *http.Request) (*models.TargetRoute,error){
	return &models.TargetRoute{Scheme: "http", Host: "paperpalprod.onrender.com", Path: "/health" },nil
}


func (proxyServiceImpl *ProxyServiceImpl) FindTargetRouteForRequest2(request *http.Request) (*models.TargetRoute,error){

	
	route,err := CheckRouteAndPath(proxyServiceImpl.ConfigManager.Get(),request.URL.Path);

	if err != nil{
		return nil,err
	}


	_ , err = getServiceForRoute(proxyServiceImpl.ConfigManager.Get(),route.Path)

	if err != nil{
		return nil,err
	}


	return &models.TargetRoute{Scheme: "http", Host: "paperpalprod.onrender.com", Path: "/health" },nil
}

func CheckRouteAndPath(cfg *models.GatewayConfig, path string) (*models.RouteConfig, error){
	for _ , route := range cfg.Routes{
		
		a := strings.TrimSuffix(path,"/**")

		if path == a {
			return &route,nil
		}
	} 

	return nil,errors.New("No Routing Config Found for Path : "+path);
}

func getServiceForRoute(cfg *models.GatewayConfig, path string) (*models.ServiceConfig,error){
	
	routeConfig,err := CheckRouteAndPath(cfg,path)

	if err != nil{
		return nil,err
	}

	serviceName := routeConfig.Service

	for _,service := range cfg.Services{

		if serviceName == service.ID {
				return &service,nil
		}

	}

	return nil,errors.New("No Targets found for URL "+path)
}