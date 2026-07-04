package middleware

import "net/http"

func Authentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//perfrom Authentiation logic and add int he context the data 
		
		next.ServeHTTP(w,r)
	})

}