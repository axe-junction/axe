package models

import "github.com/google/uuid"

type Station struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string    `json:"name"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Type      string    `json:"type"`
}

type Stations []Station

type StationRepo interface {
	GetAll() (Stations, error)
	GetByID(id string) (Stations, error)
	GetNearby(latitude, longitude float64) (Stations, error)
}
