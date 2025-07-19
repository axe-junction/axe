package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the UAA model
type User struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	GoogleID      string         `json:"google_id" gorm:"uniqueIndex;not null"`
	Email         string         `json:"email" gorm:"uniqueIndex;not null"`
	Name          string         `json:"name" gorm:"not null"`
	GivenName     string         `json:"given_name"`
	FamilyName    string         `json:"family_name"`
	Picture       string         `json:"picture"`
	Locale        string         `json:"locale"`
	VerifiedEmail bool           `json:"verified_email" gorm:"default:false"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	LastLogin     *time.Time     `json:"last_login"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserProfile contains additional user profile information
type UserProfile struct {
	ID                   uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID               uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	PhoneNumber          string         `json:"phone_number"`
	DateOfBirth          *time.Time     `json:"date_of_birth"`
	Gender               string         `json:"gender"`
	Address              string         `json:"address"`
	City                 string         `json:"city"`
	Country              string         `json:"country" gorm:"default:'Algeria'"`
	PreferredLanguage    string         `json:"preferred_language" gorm:"default:'fr'"`
	NotificationsEnabled bool           `json:"notifications_enabled" gorm:"default:true"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`

	// Foreign key relationship
	User User `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// UserSession represents user sessions for tracking
type UserSession struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	SessionID string         `json:"session_id" gorm:"uniqueIndex;not null"`
	IPAddress string         `json:"ip_address"`
	UserAgent string         `json:"user_agent"`
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Foreign key relationship
	User User `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// UserRole represents roles in the system
type UserRole struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserRoleAssignment represents user role assignments
type UserRoleAssignment struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	RoleID    uuid.UUID      `json:"role_id" gorm:"type:uuid;not null;index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Foreign key relationships
	User User     `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Role UserRole `json:"role" gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Permission represents permissions in the system
type Permission struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	Resource    string         `json:"resource" gorm:"not null"`
	Action      string         `json:"action" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// RolePermission represents role-permission mappings
type RolePermission struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RoleID       uuid.UUID      `json:"role_id" gorm:"type:uuid;not null;index"`
	PermissionID uuid.UUID      `json:"permission_id" gorm:"type:uuid;not null;index"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Foreign key relationships
	Role       UserRole   `json:"role" gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Permission Permission `json:"permission" gorm:"foreignKey:PermissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
