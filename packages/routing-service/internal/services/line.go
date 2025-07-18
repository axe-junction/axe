package services

import (
	"context"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/google/uuid"
)

type LineService struct {
	repo models.LigneRepo
}

func NewLineService(repo models.LigneRepo) *LineService {
	return &LineService{repo: repo}
}

func (s *LineService) GetAllLignes() (models.Lignes, error) {
	return s.repo.GetAll()
}

func (s *LineService) GetLigneByID(id uuid.UUID) (models.Ligne, error) {
	return s.repo.GetByID(id)
}

func (s *LineService) GetLignesByType(typee string) (models.Lignes, error) {
	return s.repo.GetByType(typee)
}

func (s *LineService) GetRouteBetweenStations(ctx context.Context, startID, endID uuid.UUID) (models.Lignes, error) {
	return s.repo.GetRouteBetweenStations(ctx, startID, endID)
}
