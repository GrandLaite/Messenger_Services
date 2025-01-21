package main

import (
	"log"
	"logging-service/internal/service"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	kafkaBroker := os.Getenv("LOG_KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}
	esAddr := os.Getenv("LOG_ES_ADDR")
	if esAddr == "" {
		esAddr = "http://localhost:9200"
	}
	esIndex := os.Getenv("LOG_ES_INDEX")
	if esIndex == "" {
		esIndex = "logs"
	}

	s := service.NewElasticService(esAddr, esIndex)
	c := service.NewKafkaConsumer(kafkaBroker, "logs", s)
	go c.Start()
	log.Println("Logging service started")
	select {}
}
