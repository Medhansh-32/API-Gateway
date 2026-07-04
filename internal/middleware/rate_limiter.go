package middleware

import "net/http"

func RateLimiter(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//perfrom RateLimit logic and add int he context the data 

		next.ServeHTTP(w,r)
	})

}