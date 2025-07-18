package models

type Ligne struct{
	ID        string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name	  string  `json:"name"`
	Stations  Stations `json:"stations" gorm:"many2many:ligne_stations;"`
	Type string  `json:"type"`
	// Distance float64 `json:"distance"` 
	StartGeo float64 `json:"start_geo"` 
	EndGeo   float64 `json:"end_geo"`  

}
type Lignes []Ligne


type LigneRepo interface{
	GetLignes() (Lignes, error)
	GetByID(id string) (Ligne, error)
	GetByType(typee string) (Lignes, error)

}