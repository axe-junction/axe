package database

import (
	"log"

	"github.com/axe-junction/axe-server/internal/models"
	"gorm.io/gorm"
)

// MigrateEnhancedTables creates all the new tables for enhanced functionality
func MigrateEnhancedTables(db *gorm.DB) error {
	log.Println("Running enhanced database migrations...")

	// Migrate all the new models
	err := db.AutoMigrate(
		&models.Review{},
		&models.Payment{},
		&models.PaymentTransaction{},
		&models.UserContributedStop{},
		&models.StopDemandRequest{},
		&models.CachedRoute{},
	)

	if err != nil {
		return err
	}

	// Add indexes for better performance
	if err := addIndexes(db); err != nil {
		log.Printf("Warning: Failed to add some indexes: %v", err)
	}

	// Update existing tables with new fields
	if err := updateExistingTables(db); err != nil {
		log.Printf("Warning: Failed to update existing tables: %v", err)
	}

	log.Println("Enhanced database migrations completed successfully")
	return nil
}

func addIndexes(db *gorm.DB) error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_reviews_station_id ON reviews(station_id);",
		"CREATE INDEX IF NOT EXISTS idx_reviews_line_id ON reviews(line_id);",
		"CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_reviews_created_at ON reviews(created_at);",

		"CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);",
		"CREATE INDEX IF NOT EXISTS idx_payments_created_at ON payments(created_at);",

		"CREATE INDEX IF NOT EXISTS idx_payment_transactions_payment_id ON payment_transactions(payment_id);",
		"CREATE INDEX IF NOT EXISTS idx_payment_transactions_status ON payment_transactions(status);",

		"CREATE INDEX IF NOT EXISTS idx_user_contributed_stops_location ON user_contributed_stops(latitude, longitude);",
		"CREATE INDEX IF NOT EXISTS idx_user_contributed_stops_status ON user_contributed_stops(status);",
		"CREATE INDEX IF NOT EXISTS idx_user_contributed_stops_demand_score ON user_contributed_stops(demand_score);",
		"CREATE INDEX IF NOT EXISTS idx_user_contributed_stops_contributor ON user_contributed_stops(contributor_id);",

		"CREATE INDEX IF NOT EXISTS idx_stop_demand_requests_stop_id ON stop_demand_requests(user_contributed_stop_id);",
		"CREATE INDEX IF NOT EXISTS idx_stop_demand_requests_user_id ON stop_demand_requests(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_stop_demand_requests_requested_at ON stop_demand_requests(requested_at);",

		"CREATE INDEX IF NOT EXISTS idx_cached_routes_hash ON cached_routes(route_hash);",
		"CREATE INDEX IF NOT EXISTS idx_cached_routes_expires_at ON cached_routes(expires_at);",
		"CREATE INDEX IF NOT EXISTS idx_cached_routes_location ON cached_routes(from_lat, from_lng, to_lat, to_lng);",
	}

	for _, index := range indexes {
		if err := db.Exec(index).Error; err != nil {
			log.Printf("Failed to create index: %s - %v", index, err)
		}
	}

	return nil
}

func updateExistingTables(db *gorm.DB) error {
	// Add new columns to existing tables
	updates := []string{
		// Update stations table with new fields
		"ALTER TABLE stations ADD COLUMN IF NOT EXISTS safety_rating DECIMAL(3,2) DEFAULT 3.0;",
		"ALTER TABLE stations ADD COLUMN IF NOT EXISTS review_count INTEGER DEFAULT 0;",
		"ALTER TABLE stations ADD COLUMN IF NOT EXISTS is_verified BOOLEAN DEFAULT true;",
		"ALTER TABLE stations ADD COLUMN IF NOT EXISTS contributor_id UUID;",
		"ALTER TABLE stations ADD COLUMN IF NOT EXISTS demand_score DECIMAL(10,2) DEFAULT 0.0;",
		"ALTER TABLE stations ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;",
		"ALTER TABLE stations ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;",

		// Update lines table with new fields (note: this assumes you have a lines table)
		// If you're using the lignes table instead, you might need different updates
		"ALTER TABLE lines ADD COLUMN IF NOT EXISTS base_fare DECIMAL(10,2) DEFAULT 0.0;",
		"ALTER TABLE lines ADD COLUMN IF NOT EXISTS fare_per_km DECIMAL(10,2) DEFAULT 0.0;",
		"ALTER TABLE lines ADD COLUMN IF NOT EXISTS payment_methods TEXT[];",
		"ALTER TABLE lines ADD COLUMN IF NOT EXISTS safety_rating DECIMAL(3,2) DEFAULT 3.0;",
		"ALTER TABLE lines ADD COLUMN IF NOT EXISTS review_count INTEGER DEFAULT 0;",
		"ALTER TABLE lines ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;",
		"ALTER TABLE lines ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;",
		"ALTER TABLE lines ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;",
	}

	for _, update := range updates {
		if err := db.Exec(update).Error; err != nil {
			log.Printf("Failed to execute update: %s - %v", update, err)
		}
	}

	return nil
}

// SeedEnhancedData seeds the database with sample data for the enhanced features
func SeedEnhancedData(db *gorm.DB) error {
	log.Println("Seeding enhanced data...")

	// Update stations with sample safety ratings and review counts
	db.Exec("UPDATE stations SET safety_rating = 3.5 + RANDOM() * 1.5 WHERE safety_rating = 0 OR safety_rating IS NULL;")
	db.Exec("UPDATE stations SET review_count = FLOOR(RANDOM() * 20) WHERE review_count = 0 OR review_count IS NULL;")

	// If you have a lines table, update it with sample payment methods and pricing
	db.Exec("UPDATE lines SET base_fare = 30.0 + RANDOM() * 20.0 WHERE base_fare = 0 OR base_fare IS NULL;")
	db.Exec("UPDATE lines SET fare_per_km = 2.0 + RANDOM() * 3.0 WHERE fare_per_km = 0 OR fare_per_km IS NULL;")
	db.Exec("UPDATE lines SET safety_rating = 3.0 + RANDOM() * 2.0 WHERE safety_rating = 0 OR safety_rating IS NULL;")
	db.Exec("UPDATE lines SET review_count = FLOOR(RANDOM() * 50) WHERE review_count = 0 OR review_count IS NULL;")

	log.Println("Enhanced data seeding completed")
	return nil
}

// CleanupExpiredData removes old cached routes and processes pending payments
func CleanupExpiredData(db *gorm.DB) error {
	// Remove expired cached routes
	db.Exec("DELETE FROM cached_routes WHERE expires_at < NOW();")

	// Process old pending payments (mark as failed after 24 hours)
	db.Exec("UPDATE payments SET status = 'failed' WHERE status = 'pending' AND created_at < NOW() - INTERVAL '24 hours';")

	return nil
}
