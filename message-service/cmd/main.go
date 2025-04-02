package main

import (
	"log"
	"net/http"
	"os"

	"message-service/internal/handlers"
	"message-service/internal/metrics"
	"message-service/internal/repository"
	"message-service/internal/service"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	port := os.Getenv("MESSAGE_SERVICE_PORT")
	if port == "" {
		port = "8083"
	}
	dbURL := os.Getenv("MESSAGE_DB_URL")
	if dbURL == "" {
		dbURL = "postgres://root:root@localhost:5432/main_db?sslmode=disable"
	}
	db, err := repository.NewDB(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	msgRepo := repository.NewMessageRepository(db)
	msgSrv := service.NewMessageService(msgRepo)
	hnd := handlers.NewMessageHandlers(msgSrv)

	r := mux.NewRouter()
	r.Use(metrics.MetricsMiddleware)

	// Эндпоинты для message-service
	r.HandleFunc("/messages/create", hnd.CreateMessageHandler).Methods("POST")
	// убрали /messages/list
	r.HandleFunc("/messages/get/{id}", hnd.GetMessageHandler).Methods("GET")
	r.HandleFunc("/messages/delete/{id}", hnd.DeleteMessageHandler).Methods("DELETE")
	r.HandleFunc("/messages/like/{id}", hnd.LikeMessageHandler).Methods("POST")
	r.HandleFunc("/messages/unlike/{id}", hnd.UnlikeMessageHandler).Methods("DELETE")
	r.HandleFunc("/messages/superlike/{id}", hnd.SuperlikeMessageHandler).Methods("POST")
	r.HandleFunc("/messages/unsuperlike/{id}", hnd.UnsuperlikeMessageHandler).Methods("DELETE")
	r.HandleFunc("/messages/conversation/{partner}", hnd.ConversationHandler).Methods("GET")
	r.HandleFunc("/messages/dialogs", hnd.DialogsHandler).Methods("GET")

	// Прометеевский эндпоинт
	r.Handle("/metrics", metrics.PrometheusHandler()).Methods("GET")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	log.Printf("Message service on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
