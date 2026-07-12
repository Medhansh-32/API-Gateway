package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/medhansh-32/api-gateway/internal/models/requests"
	"github.com/medhansh-32/api-gateway/internal/service"
)


type GateWayHandler struct{
	proxyService service.ProxyService 
	authService service.AuthenticationService
}

func (gateWayHandler *GateWayHandler) RegisteRoutes(router *http.ServeMux) {
	router.HandleFunc("/", gateWayHandler.Redirect)
	router.HandleFunc("/login", gateWayHandler.Login)
}

func NewGateWayHandler(proxyService service.ProxyService,
	authService service.AuthenticationService) (*GateWayHandler){
	return &GateWayHandler{proxyService: proxyService,authService: authService}
}

func (gateWayHandler *GateWayHandler) Redirect(w http.ResponseWriter, r *http.Request) {

    targetRoute, err := gateWayHandler.proxyService.FindTargetRouteForRequest(r)
    if err != nil {
		log.Println(err.Error())
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid Request"))
        return
    }

    targetURL := &url.URL{Scheme: "https", Host: targetRoute.Host}

    reverseProxy := &httputil.ReverseProxy{
        Rewrite: func(pr *httputil.ProxyRequest) {
            pr.SetURL(targetURL)  
            pr.Out.URL.Path = targetRoute.Path
		},
    }
	log.Println("Calling Target URL ")
    reverseProxy.ServeHTTP(w, r)
}


func (gateWayHandler *GateWayHandler) Login(w http.ResponseWriter, r *http.Request) {

	var loginRequest requests.LoginRequest

	json.NewDecoder(r.Body).Decode(&loginRequest)

	loginResponse, err := gateWayHandler.authService.Login(loginRequest)

	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse)
}
