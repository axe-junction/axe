package repo

import (
	"context"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LineRepo struct {
	db *gorm.DB
}

func NewLineRepo(db *gorm.DB) *LineRepo {
	return &LineRepo{db: db}
}

func (r *LineRepo) GetAll() (models.Lines, error) {
	var lines []models.Line
	err := r.db.Find(&lines).Error
	return lines, err
}

func (r *LineRepo) GetByID(id uuid.UUID) (models.Line, error) {
	var line models.Line
	err := r.db.Where("id = ?", id).First(&line).Error
	return line, err
}

func (r *LineRepo) GetByType(typee string) (models.Lines, error) {
	var lines []models.Line
	err := r.db.Where("type = ?", typee).Find(&lines).Error
	return lines, err
}

func (r *LineRepo) GetRouteBetweenStations(ctx context.Context, startID, endID uuid.UUID) (models.Lines, error) {
	// This can be implemented later if needed
	return []models.Line{}, nil
}

func (r *LineRepo) GetStopsForLine(lineID uuid.UUID) ([]models.Stop, error) {
	var stops []models.Stop
	err := r.db.Preload("Station").
		Where("line_id = ?", lineID).
		Order("sequence").
		Find(&stops).Error
	return stops, err
}

