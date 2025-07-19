package handlers

import (
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

type PriceResult struct {
    From  string  `json:"from"`
    To    string  `json:"to"`
    Price float64 `json:"price"`
}

type EstimateRequest struct {
    From string `json:"from"`
    To   string `json:"to"`
}

func StartKafkaEstimateWorker() error {
    consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": "localhost:9092",
        "group.id":          "estimate-service-worker",
        "auto.offset.reset": "earliest",
    })
    if err != nil {
        return fmt.Errorf("failed to create consumer: %w", err)
    }

    producer, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": "localhost:9092",
    })
    if err != nil {
        return fmt.Errorf("failed to create producer: %w", err)
    }

    err = consumer.Subscribe("estimate-request", nil)
    if err != nil {
        return fmt.Errorf("failed to subscribe: %w", err)
    }

    for {
        msg, err := consumer.ReadMessage(-1)
        if err != nil {
            log.Printf("Consumer error: %v\n", err)
            continue
        }

        var req EstimateRequest
        if err := json.Unmarshal(msg.Value, &req); err != nil {
            log.Printf(" Invalid JSON: %v\n", err)
            continue
        }
        //call the estime service 
        price := PriceResult{
            From:  req.From,
            To:    req.To,
            Price: 150 + rand.Float64()*50,
        }

        responseBytes, _ := json.Marshal([]PriceResult{price})

        var correlationID string
        for _, h := range msg.Headers {
            if h.Key == "correlation_id" {
                correlationID = string(h.Value)
                break
            }
        }

        responseMsg := kafka.Message{
            TopicPartition: kafka.TopicPartition{Topic: strPtr("estimate-response"), Partition: kafka.PartitionAny},
            Value:          responseBytes,
            Headers: []kafka.Header{
                {Key: "correlation_id", Value: []byte(correlationID)},
            },
        }

        err = producer.Produce(&responseMsg, nil)
        if err != nil {
            log.Printf(" Failed to send response: %v", err)
        } else {
            fmt.Printf("Sent estimate response for %s → %s [correlation_id=%s]\n", req.From, req.To, correlationID)
        }
    }
}

func strPtr(s string) *string {
    return &s
}
