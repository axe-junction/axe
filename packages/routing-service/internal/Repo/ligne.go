package repo

import (
	"github.com/axe-junction/axe-server/internal/models"
	"gorm.io/gorm"
)


type LigneRepo struct {
	db *gorm.DB

}
func NewLigneRepo(db *gorm.DB) *LigneRepo {
	return &LigneRepo{db: db}
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
func (repo *LigneRepo) GetNearby(latitude, longitude float64) (models.Lignes, error) {
	 radius := 20 
	var lignes models.Lignes
	for i:=radius; i>0 ; i++{
		if err := repo.db.Preload("Stations").Where("ST_DWithin(ST_MakePoint(longitude, latitude)::geography, ST_MakePoint(?, ?)::geography, ?)", longitude, latitude, radius).Find(&lignes).Error; err != nil {
			return nil, err
		}
		if len(lignes) == 0 {
			radius += 100
			continue
		}
		break
	}


	return lignes, nil
}



var _ models.LigneRepo = (*LigneRepo)(nil)