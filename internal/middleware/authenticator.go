package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/medhansh-32/api-gateway/internal/config"
	"github.com/medhansh-32/api-gateway/internal/models"
	"github.com/medhansh-32/api-gateway/internal/models/response"
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

		rWithContext, err := autheMiddleware.authenticateRequest(r, autheMiddleware.gateWayConfig)

		if err != nil {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		log.Println("Authentication Passed....")
		next.ServeHTTP(w, rWithContext)
	})

}

func (a AuthMiddleware) authenticateRequest(r *http.Request, cfg *config.ConfigManager) (*http.Request, error) {

	url := r.URL

	log.Println("Authenticating Request for URL : {}", url, " token ")

	routeConfig := getRouteConfigForURL(url, cfg.Get())

	if routeConfig == nil || routeConfig.Auth.Enabled == false {
		return r, nil
	}

	bearerToken := r.Header.Get(utils.AUTHORIZE_TOKEN)

	if bearerToken == "" {
		log.Println("Token not found in the Header")
		return nil, errors.New("Auth Token Not Found")
	}

	token := strings.TrimPrefix(bearerToken, utils.BEARER)

	log.Println(token)

	log.Println("Authenticating Request for URL : ", url, " token :", token)

	claims, err := a.authService.ValidateToken(token)

	if err != nil {
		return nil, err
	}

	log.Printf("Claims: %+v\n", claims)

	if routeConfig.Auth.Type == "Authorize" {
		contains := slices.Contains(routeConfig.Auth.Roles, claims.Role)

		if !contains {
			return nil, errors.New("Unathourized User")
		}
	}

	context := context.WithValue(r.Context(), utils.USER_INFO, claims)
	log.Println("User Authenticated user-id : ", claims.ID)
	return r.WithContext(context), nil

}

func getRouteConfigForURL(url *url.URL, cfg *models.GatewayConfig) *models.RouteConfig {
	routes := cfg.Routes

	for _, route := range routes {

		if matchURL(url, route.Path) {
			log.Println("Route Matched : ", route)
			return &route
		}

	}

	return nil

}

func matchURL(url *url.URL, pathPattern string) bool {

	path := strings.TrimSuffix(pathPattern, "/**")

	if path == getFirstSegment(url.Path) {
		return true
	}

	return false
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
