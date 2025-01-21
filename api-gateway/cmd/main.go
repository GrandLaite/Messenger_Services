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

	// ---- Маршруты БЕЗ проверки токена ----
	// Всё, что /auth/... (login/register)
	r.HandleFunc("/auth/{rest:.*}", gateway.ProxyHandler("AUTH_SERVICE_URL")).Methods("GET", "POST", "PUT", "DELETE")

	// ---- Маршруты С проверкой токена ----
	api := r.PathPrefix("/").Subrouter()
	api.Use(gateway.JWTMiddleware)
	api.HandleFunc("/users/{rest:.*}", gateway.ProxyHandler("USER_SERVICE_URL")).Methods("GET", "POST", "PUT", "DELETE")
	api.HandleFunc("/messages/{rest:.*}", gateway.ProxyHandler("MESSAGE_SERVICE_URL")).Methods("GET", "POST", "PUT", "DELETE")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("API Gateway on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
