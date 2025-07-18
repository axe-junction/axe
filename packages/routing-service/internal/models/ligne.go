package models

import (
	"context"

	"github.com/google/uuid"
)

type Ligne struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Stops    []Stop    `json:"stops" gorm:"foreignKey:RouteID"`
	Stations []Station `json:"stations" gorm:"many2many:stops;joinForeignKey:RouteID;joinReferences:StationID"`
}
type Lignes []Ligne

type LigneRepo interface {
	GetLignes() (Lignes, error)
	GetByID(id uuid.UUID) (Ligne, error)
	GetByType(typee string) (Lignes, error)
	GetRouteBetweenStations(ctx context.Context, startID, endID uuid.UUID) (Lignes, error)
}
