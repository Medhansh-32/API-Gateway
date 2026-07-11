package middleware

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//perfrom Logger logic and add int he context the data 
		log.Println("Logging Passed....")
		next.ServeHTTP(w,r)
	})

}