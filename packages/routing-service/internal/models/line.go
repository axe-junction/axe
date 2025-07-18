package models

import (
	"context"

	"github.com/google/uuid"
)

type Ligne struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name     string
	Type     string
	Agency   string    // "ETUSA", "SNTF", etc.
	Stations []Station `gorm:"many2many:ligne_stations;"`
}

type Lignes []Ligne

type LigneRepo interface {
	GetAll() (Lignes, error)
	GetByID(id uuid.UUID) (Ligne, error)
	GetByType(typee string) (Lignes, error)
	GetRouteBetweenStations(ctx context.Context, startID, endID uuid.UUID) (Lignes, error)
	GetStopsForLine(lineID uuid.UUID) ([]Stop, error)
}
