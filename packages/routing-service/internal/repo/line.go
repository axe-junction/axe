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

func NewLineRepo(db *gorm.DB) models.LineRepo {
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
	// Find lines that serve both stations
	var lines []models.Line
	err := r.db.Preload("Stations").
		Joins("JOIN stops s1 ON lines.id = s1.line_id").
		Joins("JOIN stops s2 ON lines.id = s2.line_id").
		Where("s1.station_id = ? AND s2.station_id = ?", startID, endID).
		Find(&lines).Error

	return lines, err
}

func (r *LineRepo) GetStopsForLine(lineID uuid.UUID) ([]models.Stop, error) {
	var stops []models.Stop
	err := r.db.Preload("Station").
		Where("line_id = ?", lineID).
		Order("sequence").
		Find(&stops).Error
	return stops, err
}

// Enhanced methods for safety filtering and payment methods
func (r *LineRepo) GetLinesWithSafetyFilter(minSafetyRating float64) (models.Lines, error) {
	var lines []models.Line
	err := r.db.Where("safety_rating >= ?", minSafetyRating).Find(&lines).Error
	return lines, err
}

func (r *LineRepo) UpdateLineSafetyRating(lineID uuid.UUID, rating float64) error {
	return r.db.Model(&models.Line{}).
		Where("id = ?", lineID).
		Update("safety_rating", rating).Error
}

func (r *LineRepo) GetLinesByPaymentMethod(paymentMethods []string) (models.Lines, error) {
	var lines []models.Line
	// Use PostgreSQL array overlap operator to check if any payment method matches
	err := r.db.Where("payment_methods && ?", "{"+joinStrings(paymentMethods, ",")+"}").Find(&lines).Error
	return lines, err
}

// Helper function to join strings
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for _, s := range strs[1:] {
		result += sep + s
	}
	return result
}
