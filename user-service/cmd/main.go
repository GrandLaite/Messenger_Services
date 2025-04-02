package main

import (
	"log"
	"net/http"
	"os"

	"user-service/internal/handlers"
	"user-service/internal/metrics"
	"user-service/internal/repository"
	"user-service/internal/service"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("USER_SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	dbURL := os.Getenv("USER_DB_URL")
	if dbURL == "" {
		dbURL = "postgres://root:root@localhost:5432/main_db?sslmode=disable"
	}

	db, err := repository.NewDB(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)
	srvUser := service.NewUserService(repo)
	hnd := handlers.NewUserHandlers(srvUser)

	r := mux.NewRouter()
	r.Use(metrics.MetricsMiddleware)

	// Эндпоинты для user-service
	r.HandleFunc("/users/create", hnd.CreateUserHandler).Methods("POST")
	r.HandleFunc("/users/checkpassword", hnd.CheckPasswordHandler).Methods("POST")
	r.HandleFunc("/users/search/{nickname}", hnd.SearchUserHandler).Methods("GET")

	// Эндпоинт для метрик
	r.Handle("/metrics", metrics.PrometheusHandler()).Methods("GET")

	srvHTTP := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("User service on port %s\n", port)
	log.Fatal(srvHTTP.ListenAndServe())
}
