package handlers

import (
	"log"
	"net/http"
	"time"

	"backend/config"
	"backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Error constants for consistent error messages
const (
	ErrInvalidCredentials      = "Invalid credentials"
	ErrAuthorizationRequired   = "Authorization header required"
	ErrUserNotAuthenticated    = "User not authenticated"
	ErrUserNotFound            = "User not found"
	ErrFailedToGenerateToken   = "Failed to generate token"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt int64        `json:"expiresAt"`
	User      UserResponse `json:"user"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		log.Printf("Login attempt failed: user %s not found", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidCredentials})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Login attempt failed: invalid password for user %s", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidCredentials})
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	})

	cfg := config.LoadConfig()
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Printf("Failed to generate token for user %s: %v", user.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrFailedToGenerateToken})
		return
	}

	log.Printf("User %s logged in successfully", user.Username)
	c.JSON(http.StatusOK, LoginResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt.Unix(),
		User: UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
	})
}

// Logout invalidates the token (client-side token removal + server confirmation)
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// GetCurrentUser returns the authenticated user's info
func GetCurrentUser(c *gin.Context) {
	userID := getUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrAuthorizationRequired})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		log.Printf("GetCurrentUser: user %s not found", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": ErrUserNotFound})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	})
}
