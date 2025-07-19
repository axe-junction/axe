package services

import (
	"context"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/google/uuid"
)

type LineService struct {
	repo models.LineRepo
}

func NewLineService(repo models.LineRepo) *LineService {
	return &LineService{repo: repo}
}

func (s *LineService) GetAllLines() (models.Lines, error) {
	return s.repo.GetAll()
}

func (s *LineService) GetLineByID(id uuid.UUID) (models.Line, error) {
	return s.repo.GetByID(id)
}

func (s *LineService) GetLinesByType(typee string) (models.Lines, error) {
	return s.repo.GetByType(typee)
}

func (s *LineService) GetRouteBetweenStations(ctx context.Context, startID, endID uuid.UUID) (models.Lines, error) {
	return s.repo.GetRouteBetweenStations(ctx, startID, endID)
}

// Enhanced methods
func (s *LineService) GetLinesWithSafetyFilter(minSafetyRating float64) (models.Lines, error) {
	return s.repo.GetLinesWithSafetyFilter(minSafetyRating)
}

func (s *LineService) UpdateLineSafetyRating(lineID uuid.UUID, rating float64) error {
	return s.repo.UpdateLineSafetyRating(lineID, rating)
}

func (s *LineService) GetLinesByPaymentMethod(paymentMethods []string) (models.Lines, error) {
	return s.repo.GetLinesByPaymentMethod(paymentMethods)
}