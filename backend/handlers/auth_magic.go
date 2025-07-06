package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"playslate-backend/models"
	"playslate-backend/services"
)

func HandleMagicLinkRequest(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.MagicLinkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate secure token
		token, err := generateSecureToken(32)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Create magic link with 15-minute expiration
		magicLink := models.MagicLink{
			Email:     req.Email,
			Token:     token,
			ExpiresAt: time.Now().Add(15 * time.Minute),
		}

		if err := db.Create(&magicLink).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create magic link"})
			return
		}

		// Create the magic link URL
		baseURL := os.Getenv("FRONTEND_URL")
		if baseURL == "" {
			baseURL = "http://localhost:5173"
		}
		magicLinkURL := fmt.Sprintf("%s/auth/verify?token=%s", baseURL, token)

		// Send email with magic link
		if err := services.SendMagicLinkEmail(req.Email, magicLinkURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send magic link email"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Magic link sent to your email",
			"email":   req.Email,
		})
	}
}

func HandleMagicLinkVerify(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.MagicLinkVerifyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find and validate magic link
		var magicLink models.MagicLink
		if err := db.Where("token = ? AND used_at IS NULL", req.Token).First(&magicLink).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired magic link"})
			return
		}

		// Check if token is expired
		if time.Now().After(magicLink.ExpiresAt) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Magic link has expired"})
			return
		}

		// Mark magic link as used
		now := time.Now()
		magicLink.UsedAt = &now
		if err := db.Save(&magicLink).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark magic link as used"})
			return
		}

		// Find or create user
		var user models.User
		if err := db.Where("email = ?", magicLink.Email).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create new user
				user = models.User{
					Email:    magicLink.Email,
					IsActive: true,
				}
				if err := db.Create(&user).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}
		}

		// Update last login
		now = time.Now()
		user.LastLoginAt = &now
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		// Generate JWT token
		token, err := generateToken(user.ID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, models.LoginResponse{
			Token: token,
			User:  user,
		})
	}
}

func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}