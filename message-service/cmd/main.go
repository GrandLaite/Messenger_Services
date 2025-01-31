package main

import (
	"log"
	"net/http"
	"os"

	"message-service/internal/handlers"
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
	attRepo := repository.NewAttachmentRepository(db)
	msgSrv := service.NewMessageService(msgRepo, attRepo)
	hnd := handlers.NewMessageHandlers(msgSrv)

	r := mux.NewRouter()
	r.HandleFunc("/messages", hnd.CreateMessageHandler).Methods("POST")
	r.HandleFunc("/messages", hnd.ListMessagesHandler).Methods("GET")
	r.HandleFunc("/messages/{id}", hnd.GetMessageHandler).Methods("GET")
	r.HandleFunc("/messages/update/{id}", hnd.UpdateMessageHandler).Methods("PUT")
	r.HandleFunc("/messages/delete/{id}", hnd.DeleteMessageHandler).Methods("DELETE")
	r.HandleFunc("/messages/like/{id}", hnd.LikeMessageHandler).Methods("POST")
	r.HandleFunc("/messages/superlike/{id}", hnd.SuperlikeMessageHandler).Methods("POST")
	r.HandleFunc("/messages/unlike/{id}", hnd.UnlikeMessageHandler).Methods("DELETE")
	r.HandleFunc("/messages/unsuperlike/{id}", hnd.UnsuperlikeMessageHandler).Methods("DELETE")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	log.Printf("Message service on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
