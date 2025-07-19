package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/axe-junction/axe-server/gateway/internal/config"
	"github.com/axe-junction/axe-server/gateway/internal/models"
	"github.com/axe-junction/axe-server/gateway/internal/repository"
	"github.com/axe-junction/axe-server/gateway/internal/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	oauthService *services.OAuthService
	userRepo     *repository.UserRepository
	config       *config.Config
}

func NewAuthHandler(cfg *config.Config, userRepo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		oauthService: services.NewOAuthService(cfg),
		userRepo:     userRepo,
		config:       cfg,
	}
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	state, err := generateRandomState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate state",
			"code":  "STATE_GENERATION_FAILED",
		})
		return
	}

	session := sessions.Default(c)
	session.Set("oauth_state", state)
	session.Save()

	authURL := h.oauthService.GetAuthURL(state)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
	})
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	// Verify state parameter
	session := sessions.Default(c)
	storedState := session.Get("oauth_state")
	if storedState == nil || storedState != c.Query("state") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid state parameter",
			"code":  "INVALID_STATE",
		})
		return
	}

	// Clear the state from session
	session.Delete("oauth_state")

	// Get the authorization code
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Authorization code not provided",
			"code":  "MISSING_AUTH_CODE",
		})
		return
	}

	// Exchange code for token
	token, err := h.oauthService.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange code for token",
			"code":  "TOKEN_EXCHANGE_FAILED",
		})
		return
	}

	// Get user info from Google
	userInfo, err := h.oauthService.GetUserInfo(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user info",
			"code":  "USER_INFO_FAILED",
		})
		return
	}

	// Create or update user in database
	user := &models.User{
		GoogleID:      userInfo.ID,
		Email:         userInfo.Email,
		Name:          userInfo.Name,
		GivenName:     userInfo.GivenName,
		FamilyName:    userInfo.FamilyName,
		Picture:       userInfo.Picture,
		Locale:        userInfo.Locale,
		VerifiedEmail: userInfo.VerifiedEmail,
		IsActive:      true,
	}

	err = h.userRepo.CreateOrUpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save user",
			"code":  "USER_SAVE_FAILED",
		})
		return
	}

	// Create user session
	sessionID := generateSessionID()
	userSession := &models.UserSession{
		UserID:    user.ID,
		SessionID: sessionID,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours
	}

	err = h.userRepo.CreateSession(userSession)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create session",
			"code":  "SESSION_CREATION_FAILED",
		})
		return
	}

	// Store user ID and session ID in session
	session.Set("user_id", user.ID.String())
	session.Set("session_id", sessionID)
	session.Save()

	// Return success response with user info
	c.JSON(http.StatusOK, gin.H{
		"message":    "Successfully authenticated",
		"user":       user,
		"session_id": sessionID,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	sessionID := session.Get("session_id")

	if sessionID != nil {
		// Delete session from database
		h.userRepo.DeleteSession(sessionID.(string))
	}

	// Clear session
	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
		"code":    "LOGOUT_SUCCESS",
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
			"code":  "NOT_AUTHENTICATED",
		})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID",
			"code":  "INVALID_USER_ID",
		})
		return
	}

	user, profile, err := h.userRepo.GetUserWithProfile(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user profile",
			"code":  "PROFILE_FETCH_FAILED",
		})
		return
	}

	// Get user roles
	roles, err := h.userRepo.GetUserRoles(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user roles",
			"code":  "ROLES_FETCH_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"profile": profile,
		"roles":   roles,
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
			"code":  "NOT_AUTHENTICATED",
		})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID",
			"code":  "INVALID_USER_ID",
		})
		return
	}

	var updateData struct {
		PhoneNumber          string `json:"phone_number"`
		DateOfBirth          string `json:"date_of_birth"`
		Gender               string `json:"gender"`
		Address              string `json:"address"`
		City                 string `json:"city"`
		Country              string `json:"country"`
		PreferredLanguage    string `json:"preferred_language"`
		NotificationsEnabled *bool  `json:"notifications_enabled"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
			"code":  "INVALID_JSON",
		})
		return
	}

	// Get existing profile
	_, profile, err := h.userRepo.GetUserWithProfile(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user profile",
			"code":  "PROFILE_FETCH_FAILED",
		})
		return
	}

	// Update profile fields
	if updateData.PhoneNumber != "" {
		profile.PhoneNumber = updateData.PhoneNumber
	}
	if updateData.Gender != "" {
		profile.Gender = updateData.Gender
	}
	if updateData.Address != "" {
		profile.Address = updateData.Address
	}
	if updateData.City != "" {
		profile.City = updateData.City
	}
	if updateData.Country != "" {
		profile.Country = updateData.Country
	}
	if updateData.PreferredLanguage != "" {
		profile.PreferredLanguage = updateData.PreferredLanguage
	}
	if updateData.NotificationsEnabled != nil {
		profile.NotificationsEnabled = *updateData.NotificationsEnabled
	}

	// Parse date of birth if provided
	if updateData.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", updateData.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid date format, use YYYY-MM-DD",
				"code":  "INVALID_DATE_FORMAT",
			})
			return
		}
		profile.DateOfBirth = &dob
	}

	err = h.userRepo.UpdateUserProfile(profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update profile",
			"code":  "PROFILE_UPDATE_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"profile": profile,
	})
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
