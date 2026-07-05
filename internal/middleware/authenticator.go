package middleware

import (
	"net/http"

	"github.com/medhansh-32/api-gateway/internal/service"
)


type AuthMiddleware struct {
    authService *service.AuthenticationService
}

func (autheMiddleware *AuthMiddleware) Authentication(next http.Handler) http.Handler {
	

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		
		
		next.ServeHTTP(w,r)
	})

}