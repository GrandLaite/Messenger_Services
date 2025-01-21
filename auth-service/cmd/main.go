package main

import (
	"log"
	"net/http"
	"os"

	"auth-service/internal/handlers"
	"auth-service/internal/service"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}
	secretKey := os.Getenv("AUTH_JWT_SECRET")
	if secretKey == "" {
		secretKey = "default_secret"
	}
	authSrv := service.NewAuthService(secretKey)
	authHnd := handlers.NewAuthHandlers(authSrv)
	r := mux.NewRouter()
	r.HandleFunc("/login", authHnd.LoginHandler).Methods("POST")
	r.HandleFunc("/register", authHnd.RegisterHandler).Methods("POST")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	log.Printf("Auth service on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
