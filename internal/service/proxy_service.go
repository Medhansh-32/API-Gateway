package service

import (
	"net/http"

	"github.com/medhansh-32/api-gateway/internal/models"
)


type ProxyService interface{
	FindTargetRouteForRequest(request *http.Request) (models.TargetRoute,error)
}

type ProxyServiceImpl struct{

}

func NewProxyService() (ProxyService){
	return &ProxyServiceImpl{}
}

func (proxyServiceImpl *ProxyServiceImpl) FindTargetRouteForRequest(request *http.Request) (models.TargetRoute,error){
	return models.TargetRoute{Scheme: "http", Host: "paperpalprod.onrender.com", Path: "/health" },nil
}