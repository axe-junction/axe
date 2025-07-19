package database

import (
	"fmt"

	"github.com/axe-junction/axe-server/gateway/internal/config"
	"github.com/axe-junction/axe-server/gateway/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.UserSession{},
		&models.UserRole{},
		&models.UserRoleAssignment{},
		&models.Permission{},
		&models.RolePermission{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	// Create default roles if they don't exist
	if err := createDefaultRoles(db); err != nil {
		return nil, fmt.Errorf("failed to create default roles: %w", err)
	}

	return db, nil
}

func createDefaultRoles(db *gorm.DB) error {
	roles := []models.UserRole{
		{Name: "user", Description: "Standard user role"},
		{Name: "admin", Description: "Administrator role"},
		{Name: "moderator", Description: "Moderator role"},
	}

	for _, role := range roles {
		var existingRole models.UserRole
		err := db.Where("name = ?", role.Name).First(&existingRole).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Role doesn't exist, create it
				if err := db.Create(&role).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
