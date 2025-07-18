package models

type Ligne struct{
	ID        string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name	  string  `json:"name"`
	Stations  Stations `json:"stations"`

}