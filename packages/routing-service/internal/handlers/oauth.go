package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/axe-junction/axe-server/internal/config"
	"github.com/axe-junction/axe-server/internal/models"
	"github.com/axe-junction/axe-server/internal/repo"
	"github.com/axe-junction/axe-server/internal/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
	oauthService *services.OAuthService
	userRepo     *repo.UserRepository
	config       *config.Env
}

func NewOAuthHandler(cfg *config.Env, userRepo *repo.UserRepository) *OAuthHandler {
	return &OAuthHandler{
		oauthService: services.NewOAuthService(),
		userRepo:     userRepo,
		config:       cfg,
	}
}

func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	state, err := generateRandomState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	session := sessions.Default(c)
	session.Set("oauth_state", state)
	session.Save()

	authURL := h.oauthService.GetAuthURL(state)

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	// Verify state parameter
	session := sessions.Default(c)
	storedState := session.Get("oauth_state")
	if storedState == nil || storedState != c.Query("state") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	// Clear the state from session
	session.Delete("oauth_state")
	session.Save()

	// Get the authorization code
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not provided"})
		return
	}

	// Exchange code for token
	token, err := h.oauthService.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	// Get user info from Google
	userInfo, err := h.oauthService.GetUserInfo(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
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
	}

	err = h.userRepo.CreateOrUpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	// Store user ID in session
	session.Set("user_id", user.ID)
	session.Save()

	// Return success response with user info
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully authenticated",
		"user":    user,
	})
}

func (h *OAuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *OAuthHandler) GetProfile(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, err := h.userRepo.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
