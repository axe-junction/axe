package models

import (
	"context"

	"github.com/google/uuid"
)

type Ligne struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name     string    `gorm:"column:name" json:"name"`
	Type     string    `gorm:"column:type" json:"type"`
	StartGeo float64   `gorm:"column:start_geo" json:"start_geo"`
	EndGeo   float64   `gorm:"column:end_geo" json:"end_geo"`
	Stations []Station `gorm:"many2many:ligne_stations;"`
}

type Lignes []Ligne

type LigneRepo interface {
	GetAll() (Lignes, error)
	GetByID(id uuid.UUID) (Ligne, error)
	GetByType(typee string) (Lignes, error)
	GetRouteBetweenStations(ctx context.Context, startID, endID uuid.UUID) (Lignes, error)
	GetStopsForLine(lineID uuid.UUID) (LigneStops, error)
}
