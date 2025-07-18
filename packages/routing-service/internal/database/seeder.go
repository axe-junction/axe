package database

import (
	"fmt"
	"log"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	log.Println("🌱 Starting database seeding...")

	// Clean existing data first
	log.Println("🧹 Cleaning existing data...")
	db.Exec("DELETE FROM stops")
	db.Exec("DELETE FROM lignes")
	db.Exec("DELETE FROM stations")

	// Create Stations
	stations := []models.Station{
		{ID: uuid.New(), Name: "Kouba", Latitude: 36.7322, Longitude: 3.0801, Type: "Bus"},
		{ID: uuid.New(), Name: "El Madania", Latitude: 36.7405, Longitude: 3.0709, Type: "Bus"},
		{ID: uuid.New(), Name: "Belouizdad", Latitude: 36.7486, Longitude: 3.0600, Type: "Bus"},
		{ID: uuid.New(), Name: "Didouche Mourad", Latitude: 36.7535, Longitude: 3.0582, Type: "Bus"},
		{ID: uuid.New(), Name: "Tafourah", Latitude: 36.7601, Longitude: 3.0507, Type: "Bus"},
		{ID: uuid.New(), Name: "Place des Martyrs", Latitude: 36.7710, Longitude: 3.0588, Type: "Bus"},
		{ID: uuid.New(), Name: "Bab El Oued", Latitude: 36.7833, Longitude: 3.0500, Type: "Bus"},
		{ID: uuid.New(), Name: "Rais Hamidou", Latitude: 36.7900, Longitude: 3.0410, Type: "Bus"},
	}

	log.Printf("📍 Creating %d stations...", len(stations))
	if err := db.Create(&stations).Error; err != nil {
		return fmt.Errorf("failed to seed stations: %w", err)
	}
	log.Printf("✅ Created %d stations successfully", len(stations))

	// Create Routes
	routes := []models.Ligne{
		{ID: uuid.New(), Name: "Route 1", Type: "Bus"},
		{ID: uuid.New(), Name: "Route 2", Type: "Bus"},
	}

	log.Printf("🚌 Creating %d routes...", len(routes))
	if err := db.Create(&routes).Error; err != nil {
		return fmt.Errorf("failed to seed routes: %w", err)
	}
	log.Printf("✅ Created %d routes successfully", len(routes))

	// Create Stops for Route 1 (connecting first 4 stations)
	var stops []models.Stop
	log.Println("🚏 Creating stops for Route 1...")
	for i := 0; i < 4; i++ {
		if i >= len(stations) {
			break
		}
		stop := models.Stop{
			ID:        uuid.New(),
			RouteID:   routes[0].ID,
			StationID: stations[i].ID,
			Sequence:  i + 1,
		}
		stops = append(stops, stop)
		log.Printf("   Added stop %d: %s (sequence: %d)", i+1, stations[i].Name, stop.Sequence)
	}

	// Create Stops for Route 2 (connecting stations 4-7 with proper sequence)
	log.Println("🚏 Creating stops for Route 2...")
	for i := 3; i < 7; i++ {
		if i >= len(stations) {
			break
		}
		sequenceNum := i - 3 + 1 // This will be 1, 2, 3, 4
		stop := models.Stop{
			ID:        uuid.New(),
			RouteID:   routes[1].ID,
			StationID: stations[i].ID,
			Sequence:  sequenceNum,
		}
		stops = append(stops, stop)
		log.Printf("   Added stop %d: %s (sequence: %d)", sequenceNum, stations[i].Name, stop.Sequence)
	}

	log.Printf("🚏 Creating %d stops...", len(stops))
	for i, stop := range stops {
		if err := db.Create(&stop).Error; err != nil {
			log.Printf("❌ Failed to create stop %d: %v", i+1, err)
			return fmt.Errorf("failed to seed stop %d: %w", i+1, err)
		}
		log.Printf("✅ Created stop %d: RouteID=%s, StationID=%s, Sequence=%d", i+1, stop.RouteID, stop.StationID, stop.Sequence)
	}
	log.Printf("✅ Created %d stops successfully", len(stops))

	// Debug: Let's verify the data was created correctly
	var verifyStations []models.Station
	db.Find(&verifyStations)
	log.Printf("🔍 Verification: Found %d stations in database", len(verifyStations))

	var verifyRoutes []models.Ligne
	db.Find(&verifyRoutes)
	log.Printf("🔍 Verification: Found %d routes in database", len(verifyRoutes))

	var verifyStops []models.Stop
	db.Find(&verifyStops)
	log.Printf("🔍 Verification: Found %d stops in database", len(verifyStops))

	// Show the actual stops created
	for _, stop := range verifyStops {
		log.Printf("   Stop: RouteID=%s, StationID=%s, Sequence=%d", stop.RouteID, stop.StationID, stop.Sequence)
	}

	log.Println("✅ Seeding completed successfully.")
	return nil
}
