package service

import "net/http"


type ProxyService interface{
	FindTargetRouteForRequest(request *http.Request) (string,error)
}

type ProxyServiceImpl struct{

}

func NewProxyService() (ProxyService){
	return &ProxyServiceImpl{}
}

func (proxyServiceImpl *ProxyServiceImpl) FindTargetRouteForRequest(request *http.Request) (string,error){
	return "localhost:8081",nil
}