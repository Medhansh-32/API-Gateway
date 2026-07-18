package middleware

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		body:           bytes.NewBuffer(nil),
		statusCode:     http.StatusOK,
	}
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

func Logger(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println("Request Received from Client IP : ", r.RemoteAddr, " Path : ", r.URL.Path)
		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		var obj any

		response := lrw.body.Bytes()

		if err := json.Unmarshal(response, &obj); err == nil {
			pretty, _ := json.MarshalIndent(obj, "", "  ")
			log.Printf(
				"Method=%s Path=%s Status=%d Response=\n%s",
				r.Method,
				r.URL.Path,
				lrw.statusCode,
				pretty,
			)
		} else {
	
			log.Printf(
				"Method=%s Path=%s Status=%d Response=%s",
				r.Method,
				r.URL.Path,
				lrw.statusCode,
				lrw.body.String(),
			)
		}
	})

}
