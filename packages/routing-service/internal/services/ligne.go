package services

import (
	"context"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/google/uuid"
)

type LigneService struct {
	repo models.LigneRepo
}

func NewLigneService(repo models.LigneRepo) *LigneService {
	return &LigneService{repo: repo}
}

func (s *LigneService) GetAllLignes() (models.Lignes, error) {
	return s.repo.GetLignes()
}

func (s *LigneService) GetLigneByID(id uuid.UUID) (models.Ligne, error) {
	return s.repo.GetByID(id)
}

func (s *LigneService) GetLignesByType(typee string) (models.Lignes, error) {
	return s.repo.GetByType(typee)
}

func (s *LigneService) GetRouteBetweenStations(ctx context.Context, startID, endID uuid.UUID) (models.Lignes, error) {
	return s.repo.GetRouteBetweenStations(ctx, startID, endID)
}
