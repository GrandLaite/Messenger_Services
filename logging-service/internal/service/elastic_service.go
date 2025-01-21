package service

import (
	"bytes"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticService struct {
	client *elasticsearch.Client
	index  string
}

func NewElasticService(addr, index string) *ElasticService {
	cfg := elasticsearch.Config{
		Addresses: []string{addr},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &ElasticService{client: es, index: index}
}

func (e *ElasticService) StoreLog(msg []byte) error {
	var buf bytes.Buffer
	buf.Write(msg)
	resp, err := e.client.Index(
		e.index,
		&buf,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
