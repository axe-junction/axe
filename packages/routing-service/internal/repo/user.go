package repo

import (
	"github.com/axe-junction/axe-server/internal/models"
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
			return r.db.Create(user).Error
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

func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
