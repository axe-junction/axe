package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/axe-junction/axe-server/internal/repo"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type PlaceChange struct {
	ID     string                 `json:"id"`
	Local  map[string]interface{} `json:"local"`
	Remote map[string]interface{} `json:"remote"`
}

type ConflictPayload struct {
	Type    string        `json:"type"`
	Payload []PlaceChange `json:"payload"`
}

type RabitConsumer struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   amqp091.Queue
	repo    *repo.StationRepo
}

func NewRabitConsumer(db *gorm.DB) *RabitConsumer {
	stationRepo := repo.NewStationRepo(db)
	return &RabitConsumer{
		repo: stationRepo,
	}
}

func (r *RabitConsumer) Consume() error {
	var err error
	r.conn, err = amqp091.Dial("amqp://guest:guest@localhost/")
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer func() {
		if err := r.conn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v\n", err)
		}
	}()

	r.channel, err = r.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	r.queue, err = r.channel.QueueDeclare(
		"placetopic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := r.channel.Consume(
		r.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			var conflict ConflictPayload
			err := json.Unmarshal(msg.Body, &conflict)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v\n", err)
				continue
			}

			if conflict.Type == "conflict" {
				for _, change := range conflict.Payload {
					// Helper function to safely extract values from interface{}
					getStringValue := func(key string) string {
						if val, ok := change.Local[key]; ok && val != nil {
							if str, ok := val.(string); ok {
								return str
							}
						}
						return ""
					}

					getFloat64Value := func(key string) float64 {
						if val, ok := change.Local[key]; ok && val != nil {
							if f, ok := val.(float64); ok {
								return f
							}
						}
						return 0.0
					}

					station := &models.Station{
						ID:        uuid.MustParse(change.ID),
						Name:      getStringValue("name"),
						Latitude:  getFloat64Value("latitude"),
						Longitude: getFloat64Value("longitude"),
						Type:      getStringValue("type"),
					}

					if err := r.repo.Update(station); err != nil {
						log.Printf("Failed to update station %s: %v\n", change.ID, err)
					} else {
						log.Printf("Successfully updated station %s\n", change.ID)
					}
				}
			} else {
				log.Printf("ℹ Unknown message type: %s\n", conflict.Type)
			}
		}
	}()

	log.Println(" Waiting for messages. Press Ctrl+C to exit.")
	<-make(chan struct{})
	return nil
}
