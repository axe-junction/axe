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

func NewLineRepo(db *gorm.DB) models.LigneRepo {
	return &LineRepo{db: db}
}

func (r *LineRepo) GetAll() (models.Lignes, error) {
	var lines []models.Ligne
	err := r.db.Find(&lines).Error
	return lines, err
}

func (r *LineRepo) GetByID(id uuid.UUID) (models.Ligne, error) {
	var line models.Ligne
	err := r.db.Where("id = ?", id).First(&line).Error
	return line, err
}

func (r *LineRepo) GetByType(typee string) (models.Lignes, error) {
	var lines []models.Ligne
	err := r.db.Where("type = ?", typee).Find(&lines).Error
	return lines, err
}

func (r *LineRepo) GetRouteBetweenStations(ctx context.Context, startID, endID uuid.UUID) (models.Lignes, error) {
	// Find lines that serve both stations
	var lines []models.Ligne
	err := r.db.Preload("Stations").
		Joins("JOIN stops s1 ON lignes.id = s1.line_id").
		Joins("JOIN stops s2 ON lignes.id = s2.line_id").
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

// LegacyLineRepoWrapper provides compatibility with the old LineRepo interface
type LegacyLineRepoWrapper struct {
	repo models.LigneRepo
}

func NewLegacyLineRepoWrapper(db *gorm.DB) models.LineRepo {
	return &LegacyLineRepoWrapper{
		repo: NewLineRepo(db),
	}
}

func (r *LegacyLineRepoWrapper) GetAll() (models.Lines, error) {
	lignes, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// Convert Lignes to Lines
	lines := make(models.Lines, len(lignes))
	for i, ligne := range lignes {
		lines[i] = models.Line{
			ID:     ligne.ID,
			Name:   ligne.Name,
			Type:   ligne.Type,
			Agency: ligne.Agency,
		}
	}
	return lines, nil
}

func (r *LegacyLineRepoWrapper) GetByID(id uuid.UUID) (models.Line, error) {
	ligne, err := r.repo.GetByID(id)
	if err != nil {
		return models.Line{}, err
	}

	return models.Line{
		ID:     ligne.ID,
		Name:   ligne.Name,
		Type:   ligne.Type,
		Agency: ligne.Agency,
	}, nil
}

func (r *LegacyLineRepoWrapper) GetByType(typee string) (models.Lines, error) {
	lignes, err := r.repo.GetByType(typee)
	if err != nil {
		return nil, err
	}

	// Convert Lignes to Lines
	lines := make(models.Lines, len(lignes))
	for i, ligne := range lignes {
		lines[i] = models.Line{
			ID:     ligne.ID,
			Name:   ligne.Name,
			Type:   ligne.Type,
			Agency: ligne.Agency,
		}
	}
	return lines, nil
}

func (r *LegacyLineRepoWrapper) GetRouteBetweenStations(ctx context.Context, startID, endID uuid.UUID) (models.Lines, error) {
	lignes, err := r.repo.GetRouteBetweenStations(ctx, startID, endID)
	if err != nil {
		return nil, err
	}

	// Convert Lignes to Lines
	lines := make(models.Lines, len(lignes))
	for i, ligne := range lignes {
		lines[i] = models.Line{
			ID:     ligne.ID,
			Name:   ligne.Name,
			Type:   ligne.Type,
			Agency: ligne.Agency,
		}
	}
	return lines, nil
}

func (r *LegacyLineRepoWrapper) GetStopsForLine(lineID uuid.UUID) ([]models.Stop, error) {
	return r.repo.GetStopsForLine(lineID)
}
