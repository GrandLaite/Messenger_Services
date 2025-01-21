package main

import (
	"notification-service/internal/service"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	kafkaBroker := os.Getenv("NOTIF_KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}
	smtpHost := os.Getenv("NOTIF_SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "localhost:25"
	}
	agg := service.NewAggregator(smtpHost, 20*time.Minute)
	nc := service.NewKafkaConsumer(kafkaBroker, "message.events", agg)
	go nc.Start()
	select {}
}
