package router

import (
	checkHealth "blog-system"
	"github.com/gorilla/mux"
	"net/http"
)

func Router(r *mux.Router) {
	r.HandleFunc("/",
		checkHealth.HealthCheck,
	).Methods(http.MethodGet)

	r.Use(accessControlMiddleware)
}

func accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
