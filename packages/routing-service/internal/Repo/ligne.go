package repo

import (
	"github.com/axe-junction/axe-server/internal/models"
	"gorm.io/gorm"
)


type LigneRepo struct {
	db *gorm.DB

}

func (repo *LigneRepo) GetLignes() (models.Lignes, error) {
	var lignes models.Lignes
	if err := repo.db.Preload("Stations").Find(&lignes).Error; err != nil {
		return nil, err
	}
	return lignes, nil
}

func (repo *LigneRepo) GetByID(id string) (models.Ligne, error) {
	var ligne models.Ligne
	if err := repo.db.Preload("Stations").First(&ligne, "id = ?", id).Error; err != nil {
		return models.Ligne{}, err
	}
	return ligne, nil
}

func (repo *LigneRepo) GetByType(typee string) (models.Lignes, error) {
	var lignes models.Lignes
	if err := repo.db.Preload("Stations").Where("type = ?", typee).Find(&lignes).Error; err != nil {
		return nil, err
	}
	return lignes, nil
}



// var _ models.LigneRepo = (*LigneRepo)(nil)