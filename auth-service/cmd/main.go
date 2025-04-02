package main

import (
	"log"
	"net/http"
	"os"

	"auth-service/internal/handlers"
	"auth-service/internal/metrics"
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

	// Инициализация сервиса авторизации
	authSrv := service.NewAuthService(secretKey)

	// Инициализация хендлеров
	authHnd := handlers.NewAuthHandlers(authSrv)

	// Роутер
	r := mux.NewRouter()
	r.Use(metrics.MetricsMiddleware)

	// Эндпоинты для auth-service
	r.HandleFunc("/auth/register", authHnd.RegisterHandler).Methods("POST")
	r.HandleFunc("/auth/login", authHnd.LoginHandler).Methods("POST")

	// Эндпоинт для метрик
	r.Handle("/metrics", metrics.PrometheusHandler()).Methods("GET")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("Auth service on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
