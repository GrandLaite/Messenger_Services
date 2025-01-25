package main

import (
	"log"
	"net/http"
	"os"

	"api-gateway/internal/gateway"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()

	// Прокси-маршруты
	r.HandleFunc("/auth/{rest:.*}", gateway.ProxyHandler("AUTH_SERVICE_URL")).Methods("GET", "POST", "PUT", "DELETE")
	r.HandleFunc("/users/{rest:.*}", gateway.ProxyHandler("USER_SERVICE_URL")).Methods("GET", "POST", "PUT", "DELETE")
	r.HandleFunc("/messages/{rest:.*}", gateway.ProxyHandler("MESSAGE_SERVICE_URL")).Methods("GET", "POST", "PUT", "DELETE")
	r.HandleFunc("/notifications/{rest:.*}", gateway.ProxyHandler("NOTIFICATION_SERVICE_URL")).Methods("GET", "POST", "PUT", "DELETE")
	r.HandleFunc("/caching/{rest:.*}", gateway.ProxyHandler("CACHING_SERVICE_URL")).Methods("GET", "POST", "PUT", "DELETE")

	// Метрики
	r.HandleFunc("/metrics", gateway.MetricsHandler).Methods("GET")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("API Gateway running on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
