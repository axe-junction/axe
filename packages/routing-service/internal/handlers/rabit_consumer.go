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
	repo *repo.StationRepo
}

func NewRabitConsumer() *RabitConsumer {
	repo := repo.NewStationRepo(&gorm.DB{}) // Replace nil with actual DB connection if needed
	return &RabitConsumer{
		repo: repo,
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
					r.repo.Update(&models.Station{
	                     ID: uuid.MustParse(change.ID),	
						Name: change.Local["name"].(string),
						Latitude: change.Local["latitude"].(float64),
						Longitude: change.Local["longitude"].(float64),
						// Address: change.Local["address"].(),
						Address: change.Local["address"].(models.Station),
					},
					)	
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
