package database

import (
	"log"

	"gorm.io/gorm"
)

// MigrateLigneToLine migrates the old lignes table structure to the new lines table
func MigrateLigneToLine(db *gorm.DB) error {
	log.Println("Migrating lignes table to lines table...")

	// Check if lignes table exists
	var tableExists bool
	err := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'lignes')").Scan(&tableExists).Error
	if err != nil {
		return err
	}

	if !tableExists {
		log.Println("Lignes table does not exist, skipping migration")
		return nil
	}

	// Begin transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Step 1: Create new lines table with enhanced structure
	createLinesTable := `
	CREATE TABLE IF NOT EXISTS lines (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(255) NOT NULL,
		type VARCHAR(100) NOT NULL,
		agency VARCHAR(255) DEFAULT '',
		base_fare DECIMAL(10,2) DEFAULT 0.0,
		fare_per_km DECIMAL(10,2) DEFAULT 0.0,
		payment_methods TEXT[],
		safety_rating DECIMAL(3,2) DEFAULT 3.0,
		review_count INTEGER DEFAULT 0,
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if err := tx.Exec(createLinesTable).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Step 2: Migrate data from lignes to lines
	migrateData := `
	INSERT INTO lines (id, name, type, agency, created_at, updated_at)
	SELECT 
		id, 
		name, 
		type, 
		'' as agency,
		COALESCE(created_at, CURRENT_TIMESTAMP) as created_at,
		COALESCE(updated_at, CURRENT_TIMESTAMP) as updated_at
	FROM lignes
	ON CONFLICT (id) DO NOTHING;`

	if err := tx.Exec(migrateData).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Step 3: Update stops table foreign key column names
	// Update column name from route_id to line_id if needed
	updateStopsColumn := `
	DO $$
	BEGIN
		IF EXISTS (SELECT 1 FROM information_schema.columns 
				   WHERE table_name = 'stops' AND column_name = 'route_id') THEN
			-- Rename route_id to line_id
			ALTER TABLE stops RENAME COLUMN route_id TO line_id;
		END IF;
	END $$;`

	if err := tx.Exec(updateStopsColumn).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Step 4: Update many-to-many table name from ligne_stations to line_stations
	updateManyToManyTable := `
	DO $$
	BEGIN
		IF EXISTS (SELECT 1 FROM information_schema.tables 
				   WHERE table_name = 'ligne_stations') THEN
			-- Rename table
			ALTER TABLE ligne_stations RENAME TO line_stations;
		END IF;
	END $$;`

	if err := tx.Exec(updateManyToManyTable).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Step 5: Update foreign key column names in line_stations table
	updateLineStationsColumns := `
	DO $$
	BEGIN
		-- Update ligne_id to line_id if exists
		IF EXISTS (SELECT 1 FROM information_schema.columns 
				   WHERE table_name = 'line_stations' AND column_name = 'ligne_id') THEN
			ALTER TABLE line_stations RENAME COLUMN ligne_id TO line_id;
		END IF;
	END $$;`

	if err := tx.Exec(updateLineStationsColumns).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Step 6: Drop the old lignes table after successful migration
	if err := tx.Exec("DROP TABLE IF EXISTS lignes CASCADE").Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	log.Println("Successfully migrated lignes table to lines table")
	return nil
}
