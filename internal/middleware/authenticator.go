package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/medhansh-32/api-gateway/internal/config"
	"github.com/medhansh-32/api-gateway/internal/models"
	"github.com/medhansh-32/api-gateway/internal/service"
	"github.com/medhansh-32/api-gateway/internal/utils"
)

type AuthMiddleware struct {
	authService   *service.AuthenticationService
	gateWayConfig *config.ConfigManager
}

func NewAuthMiddleWare(authService *service.AuthenticationService, gateWayConfig *config.ConfigManager) AuthMiddleware {
	return AuthMiddleware{authService: authService, gateWayConfig: gateWayConfig}
}

func (autheMiddleware *AuthMiddleware) Authentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Authentication Passed....")
		rWithContext , err := autheMiddleware.authenticateRequest(r, autheMiddleware.gateWayConfig)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		next.ServeHTTP(w, rWithContext)
	})

}

func (a AuthMiddleware) authenticateRequest(r *http.Request, cfg *config.ConfigManager) (*http.Request, error) {

	url := r.URL

	log.Println("Authenticating Request for URL : {}", url, " token ")

	authEnbaled := checkAuthEnabledForURL(url, cfg.Get())

	if authEnbaled == false {
		return nil, nil
	}

	bearerToken := r.Header.Get(utils.AUTHORIZE_TOKEN)

	if bearerToken == "" {
		log.Println("Token not found in the Header")
		return nil, errors.New("Auth Token Not Found")
	}

	token := strings.TrimPrefix(bearerToken,"bearer ")

	log.Println("Authenticating Request for URL : ", url, " token : ", token)


	claims, err:= a.authService.ValidateToken(token)

	if err == nil{
		context := context.WithValue(r.Context(), utils.USER_INFO, claims)
		log.Println("User Authenticated user-id : ",claims.ID)
		return r.WithContext(context), nil
	}

	return nil, err

}

func checkAuthEnabledForURL(url *url.URL, cfg *models.GatewayConfig) bool {
	routes := cfg.Routes

	for _, route := range routes {
	
		if matchURL(url,route.Path){
		return route.Auth.Enabled
	  }
	
	}

	return false

}

func matchURL(url *url.URL, pathPattern string) (bool) {

	path := strings.TrimSuffix(pathPattern,"/**")

	if path == url.Path{
		return true
	}

	return false
}