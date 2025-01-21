package service

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type kafkaConsumer struct {
	broker string
	topic  string
	store  *ElasticService
}

func NewKafkaConsumer(broker, topic string, st *ElasticService) *kafkaConsumer {
	return &kafkaConsumer{broker: broker, topic: topic, store: st}
}

func (kc *kafkaConsumer) Start() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kc.broker},
		Topic:   kc.topic,
		GroupID: "logging-service-group",
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		err = kc.store.StoreLog(m.Value)
		if err != nil {
			log.Println(err)
		}
	}
}
