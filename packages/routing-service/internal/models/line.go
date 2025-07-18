package models

import (
	"context"

	"github.com/google/uuid"
)

type Ligne struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name   string
	Type   string
	Agency string // "ETUSA", "SNTF", etc.
}

type Lignes []Ligne

type LigneRepo interface {
	GetAll() (Lines, error)
	GetByID(id uuid.UUID) (Line, error)
	GetByType(typee string) (Lines, error)
	GetRouteBetweenStations(ctx context.Context, startID, endID uuid.UUID) (Lines, error)
}
