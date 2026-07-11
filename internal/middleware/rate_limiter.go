package middleware

import (
	"log"
	"net/http"
)

func RateLimiter(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//perfrom RateLimit logic and add int he context the data 

		log.Println("Rate Limiter Passed....")
		next.ServeHTTP(w,r)
	})

}