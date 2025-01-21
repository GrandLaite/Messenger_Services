package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type kafkaConsumer struct {
	broker string
	topic  string
	agg    *Aggregator
}

func NewKafkaConsumer(broker, topic string, agg *Aggregator) *kafkaConsumer {
	return &kafkaConsumer{broker: broker, topic: topic, agg: agg}
}

func (kc *kafkaConsumer) Start() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kc.broker},
		Topic:   kc.topic,
		GroupID: "notif-service-group",
	})
	defer r.Close()
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		var evt struct {
			Sender    string `json:"sender"`
			Recipient string `json:"recipient"`
		}
		err = json.Unmarshal(m.Value, &evt)
		if err != nil {
			continue
		}
		if evt.Sender != "" && evt.Recipient != "" {
			kc.agg.AddMessage(evt.Sender, evt.Recipient)
		}
	}
}
