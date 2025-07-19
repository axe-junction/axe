package repository

import (
	"time"

	"github.com/axe-junction/axe-server/gateway/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateOrUpdateUser(user *models.User) error {
	var existingUser models.User
	result := r.db.Where("google_id = ?", user.GoogleID).First(&existingUser)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new user
			user.CreatedAt = time.Now()
			user.UpdatedAt = time.Now()
			if err := r.db.Create(user).Error; err != nil {
				return err
			}

			// Create default user profile
			profile := &models.UserProfile{
				UserID:               user.ID,
				PreferredLanguage:    "fr",
				NotificationsEnabled: true,
				Country:              "Algeria",
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			}
			if err := r.db.Create(profile).Error; err != nil {
				return err
			}

			// Assign default user role
			return r.AssignDefaultRole(user.ID)
		}
		return result.Error
	}

	// Update existing user
	existingUser.Email = user.Email
	existingUser.Name = user.Name
	existingUser.GivenName = user.GivenName
	existingUser.FamilyName = user.FamilyName
	existingUser.Picture = user.Picture
	existingUser.Locale = user.Locale
	existingUser.VerifiedEmail = user.VerifiedEmail
	now := time.Now()
	existingUser.LastLogin = &now
	existingUser.UpdatedAt = time.Now()

	*user = existingUser
	return r.db.Save(&existingUser).Error
}

func (r *UserRepository) GetUserByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("google_id = ?", googleID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserWithProfile(id uuid.UUID) (*models.User, *models.UserProfile, error) {
	var user models.User
	var profile models.UserProfile

	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, nil, err
	}

	err = r.db.Where("user_id = ?", id).First(&profile).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, nil, err
	}

	return &user, &profile, nil
}

func (r *UserRepository) UpdateUserProfile(profile *models.UserProfile) error {
	profile.UpdatedAt = time.Now()
	return r.db.Save(profile).Error
}

func (r *UserRepository) CreateSession(session *models.UserSession) error {
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()
	return r.db.Create(session).Error
}

func (r *UserRepository) GetSession(sessionID string) (*models.UserSession, error) {
	var session models.UserSession
	err := r.db.Where("session_id = ? AND expires_at > ?", sessionID, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *UserRepository) DeleteSession(sessionID string) error {
	return r.db.Where("session_id = ?", sessionID).Delete(&models.UserSession{}).Error
}

func (r *UserRepository) AssignDefaultRole(userID uuid.UUID) error {
	// Get the default "user" role
	var role models.UserRole
	err := r.db.Where("name = ?", "user").First(&role).Error
	if err != nil {
		return err
	}

	// Check if the user already has this role
	var existingAssignment models.UserRoleAssignment
	err = r.db.Where("user_id = ? AND role_id = ?", userID, role.ID).First(&existingAssignment).Error
	if err == nil {
		return nil // Role already assigned
	}

	// Create role assignment
	assignment := &models.UserRoleAssignment{
		UserID:    userID,
		RoleID:    role.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return r.db.Create(assignment).Error
}

func (r *UserRepository) GetUserRoles(userID uuid.UUID) ([]models.UserRole, error) {
	var roles []models.UserRole
	err := r.db.Table("user_roles").
		Joins("JOIN user_role_assignments ON user_roles.id = user_role_assignments.role_id").
		Where("user_role_assignments.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *UserRepository) HasPermission(userID uuid.UUID, resource, action string) (bool, error) {
	var count int64
	err := r.db.Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_role_assignments ON role_permissions.role_id = user_role_assignments.role_id").
		Where("user_role_assignments.user_id = ? AND permissions.resource = ? AND permissions.action = ?", userID, resource, action).
		Count(&count).Error

	return count > 0, err
}
