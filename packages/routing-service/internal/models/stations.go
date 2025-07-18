package models

type stations struct{
	ID        string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name      string `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Type 	string `json:"type"`
}

type Stations []stations

type StationRepo interface{
	GetAll() (Stations, error)
	GetByID(id string) (stations, error)
}