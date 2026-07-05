package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/medhansh-32/api-gateway/internal/service"
)


type GateWayHandler struct{
	proxyService service.ProxyService 
}

func (gateWayHandler *GateWayHandler) RegisteRoutes(router *http.ServeMux) {
	router.HandleFunc("/", gateWayHandler.Redirect)
}

func NewGateWayHandler(proxyService service.ProxyService) (*GateWayHandler){
	return &GateWayHandler{proxyService: proxyService}
}

func (gateWayHandler *GateWayHandler) Redirect(r http.ResponseWriter, w *http.Request) {

	target, err := gateWayHandler.proxyService.FindTargetRouteForRequest(w)

	if err != nil {
		r.WriteHeader(http.StatusBadRequest)
		r.Write([]byte("Invalid Request"))
		return
	}

	targetURL := url.URL{Scheme: "http",Host: target}

	reverseProxy := httputil.NewSingleHostReverseProxy(&targetURL)

	reverseProxy.ServeHTTP(r,w)
}
