package main

import (
	"log"
	"net/http"
	"os"

	"caching-service/internal/handlers"
	"caching-service/internal/repository"
	"caching-service/internal/service"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	port := os.Getenv("CACHE_SERVICE_PORT")
	if port == "" {
		port = "8084"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	repo, err := repository.NewRedisRepository(redisAddr)
	if err != nil {
		log.Fatal(err)
	}
	srv := service.NewCacheService(repo)
	hnd := handlers.NewCacheHandlers(srv)

	r := mux.NewRouter()
	r.HandleFunc("/cache", hnd.SetHandler).Methods("POST")
	r.HandleFunc("/cache", hnd.GetHandler).Methods("GET")

	s := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	log.Printf("Caching service on port %s\n", port)
	log.Fatal(s.ListenAndServe())
}
