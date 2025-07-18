package database

import (
	"math"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func calculateHaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Earth radius in meters
	φ1 := lat1 * math.Pi / 180
	φ2 := lat2 * math.Pi / 180
	Δφ := (lat2 - lat1) * math.Pi / 180
	Δλ := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func createStations() []models.Station {
	return []models.Station{
		{ID: uuid.New(), Name: "Place des Martyrs", Latitude: 36.7764, Longitude: 3.0585, Type: "tram"},
		{ID: uuid.New(), Name: "Khelifa Boukhalfa", Latitude: 36.7703, Longitude: 3.0542, Type: "tram"},
		{ID: uuid.New(), Name: "Place du 1er Mai", Latitude: 36.7636, Longitude: 3.0501, Type: "tram"},
		{ID: uuid.New(), Name: "Les Ateliers", Latitude: 36.7568, Longitude: 3.0460, Type: "tram"},
		{ID: uuid.New(), Name: "Cité Mokhtar Zerhouni", Latitude: 36.7500, Longitude: 3.0419, Type: "tram"},

		{ID: uuid.New(), Name: "Tafourah", Latitude: 36.7740, Longitude: 3.0600, Type: "metro"},
		{ID: uuid.New(), Name: "Khelifa Boukhalfa (M)", Latitude: 36.7703, Longitude: 3.0542, Type: "metro"},
		{ID: uuid.New(), Name: "1 Mai (M)", Latitude: 36.7636, Longitude: 3.0501, Type: "metro"},

		{ID: uuid.New(), Name: "Gare d'Alger", Latitude: 36.7840, Longitude: 3.0560, Type: "bus"},
		{ID: uuid.New(), Name: "Bab El Oued", Latitude: 36.7890, Longitude: 3.0500, Type: "bus"},
		{ID: uuid.New(), Name: "Université de Bab Ezzouar", Latitude: 36.7128, Longitude: 3.1825, Type: "train"},
	}
}

func Seed(db *gorm.DB) error {
	stations := createStations()
	if err := db.Create(&stations).Error; err != nil {
		return err
	}

	lines := []models.Line{
		{ID: uuid.New(), Name: "Tramway T1", Type: "tram"},
		{ID: uuid.New(), Name: "Metro Line 1", Type: "metro"},
		{ID: uuid.New(), Name: "Bus 100", Type: "bus"},
	}

	if err := db.Create(&lines).Error; err != nil {
		return err
	}

	createStops(db, stations, lines)

	createTransfers(db, stations)

	return nil
}

func createStops(db *gorm.DB, stations []models.Station, lines []models.Line) {
	// Example: Tramway T1 route
	tramStations := []string{
		"Place des Martyrs", "Khelifa Boukhalfa", "Place du 1er Mai",
		"Les Ateliers", "Cité Mokhtar Zerhouni",
	}

	for i, name := range tramStations {
		for _, station := range stations {
			if station.Name == name {
				stop := models.Stop{
					ID:        uuid.New(),
					LineID:    lines[0].ID,
					StationID: station.ID,
					Sequence:  i + 1,
				}
				db.Create(&stop)
				break
			}
		}
	}

	// Add similar for other lines...
}

func createTransfers(db *gorm.DB, stations []models.Station) {
	// Create transfers between key interchange stations
	interchanges := map[string][]string{
		"Place des Martyrs": {"Khelifa Boukhalfa", "Tafourah"},
		"Tafourah":          {"Place des Martyrs", "Khelifa Boukhalfa"},
	}

	for fromName, toNames := range interchanges {
		var fromStation models.Station
		for _, s := range stations {
			if s.Name == fromName {
				fromStation = s
				break
			}
		}

		for _, toName := range toNames {
			var toStation models.Station
			for _, s := range stations {
				if s.Name == toName {
					toStation = s
					break
				}
			}

			if fromStation.ID != uuid.Nil && toStation.ID != uuid.Nil {
				// Calculate walking distance
				distance := calculateHaversineDistance(
					fromStation.Latitude, fromStation.Longitude,
					toStation.Latitude, toStation.Longitude,
				)

				// Create transfer in both directions
				db.Create(&models.Transfer{
					FromStationID:   fromStation.ID,
					ToStationID:     toStation.ID,
					WalkingDistance: distance,
					WalkingTime:     distance / 1.4, // 5km/h walking speed
				})

				db.Create(&models.Transfer{
					FromStationID:   toStation.ID,
					ToStationID:     fromStation.ID,
					WalkingDistance: distance,
					WalkingTime:     distance / 1.4,
				})
			}
		}
	}
}
