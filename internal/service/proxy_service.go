package service

import (
	"errors"
	"log"
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
	WrrManager    *WRRManager
}

func NewProxyService(Cfg *config.ConfigManager) (ProxyService){
	return &ProxyServiceImpl{ConfigManager: Cfg,
		WrrManager: &WRRManager{services: make(map[string]*wrrState)},
	}
}


func (proxyServiceImpl *ProxyServiceImpl) FindTargetRouteForRequest(request *http.Request) (*models.TargetRoute,error){

	
	route,err := CheckRouteAndPath(proxyServiceImpl.ConfigManager.Get(),request.URL.Path);

	if err != nil{
		return nil,err
	}


	service , err := getServiceForRoute(proxyServiceImpl.ConfigManager.Get(),route.Path)

	if err != nil{
		return nil,err
	}

	target, available :=proxyServiceImpl.WrrManager.GetTarget(service.ID,service.Targets) 

	if !available{
		return nil,errors.New("No Target Available for path : "+request.URL.Path)
	}

	log.Println("Target Route Found for Service : "+ route.Service + " target : "+target.URL+" Weight : ",target.Weight)

	return &models.TargetRoute{Host: target.URL, Path: stripFirstSegment(request.URL.Path) },nil
}

func CheckRouteAndPath(cfg *models.GatewayConfig, path string) (*models.RouteConfig, error){
	for _ , route := range cfg.Routes{
		
		baseRoute := strings.TrimSuffix(route.Path,"/**")

		if getFirstSegment(path) == baseRoute {
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

func stripFirstSegment(path string) string {
	path = strings.TrimPrefix(path, "/")

	parts := strings.SplitN(path, "/", 2)

	if len(parts) == 1 {
		return "/"
	}

	return "/" + parts[1]
}

func getFirstSegment(path string) string {
	path = strings.TrimPrefix(path, "/")

	parts := strings.SplitN(path, "/", 2)

	if len(parts) == 1 {
		return "/"
	}
	log.Println("Orignal Path : " + "/" + parts[0])
	return "/" + parts[0]
}