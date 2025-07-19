package models

import "github.com/google/uuid"

// LineStop represents a stop on a specific line
type LineStop struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	LineID    uuid.UUID `gorm:"column:line_id" json:"line_id"`
	StationID uuid.UUID `gorm:"column:station_id" json:"station_id"`
	Sequence  int       `gorm:"column:sequence" json:"sequence"`
	Line      Line      `json:"line" gorm:"foreignKey:LineID"`
	Station   Station   `json:"station" gorm:"foreignKey:StationID"`
}
type LineStops []LineStop

// TableName sets the table name for LineStop
func (LineStop) TableName() string {
	return "stops"
}
