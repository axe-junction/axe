package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
)

type Price struct {
    kafkaProducer *kafka.Producer
    kafkaConsumer *kafka.Consumer
}
type PriceQuery struct {
    From string `json:"from"`
    To   string `json:"to"`
}

type PriceResult struct {
    From  string  `json:"from"`
    To    string  `json:"to"`
    Price float64 `json:"price"`
}

func NewPriceHandler() (*Price, error) {
    producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
    if err != nil {
        return nil, err
    }

    consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": "localhost:9092",
        "group.id":          "price-service",
        "auto.offset.reset": "earliest",
    })
    if err != nil {
        return nil, err
    }

    err = consumer.Subscribe("price-response", nil)
    if err != nil {
        return nil, err
    }

    return &Price{
        kafkaProducer: producer,
        kafkaConsumer: consumer,
    }, nil
}
func (p *Price) RequestPrice(ctx context.Context, from string, to string) ([]PriceResult, error) {
    correlationID := uuid.New().String()

    query := map[string]interface{}{
        "from": from,
        "to":   to,
    }

    queryBytes, _ := json.Marshal(query)

    topic := "price-request"
    p.kafkaProducer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Value:          queryBytes,
        Headers: []kafka.Header{
            {Key: "correlation_id", Value: []byte(correlationID)},
        },
    }, nil)


    for {
        msg, err := p.kafkaConsumer.ReadMessage(10 * time.Second)
        if err != nil {
            return nil, fmt.Errorf("timeout waiting for price response")
        }

        var cid string
        for _, h := range msg.Headers {
            if h.Key == "correlation_id" {
                cid = string(h.Value)
            }
        }

        if cid != correlationID {
            continue 
        }

        var prices []PriceResult
        err = json.Unmarshal(msg.Value, &prices)
        if err != nil {
            return nil, fmt.Errorf("failed to decode response: %w", err)
        }

        return prices, nil
    }
}
