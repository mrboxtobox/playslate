package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"gorm.io/gorm"

	"playslate-backend/models"
)

func HandleGoogleAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.GoogleAuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verify Google ID token
		userInfo, err := verifyGoogleIDToken(req.IdToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Google ID token"})
			return
		}

		// Find or create user
		var user models.User
		if err := db.Where("google_id = ? OR email = ?", userInfo.Id, userInfo.Email).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create new user
				user = models.User{
					Email:     userInfo.Email,
					FirstName: userInfo.GivenName,
					LastName:  userInfo.FamilyName,
					Avatar:    userInfo.Picture,
					GoogleID:  userInfo.Id,
					IsActive:  true,
				}
				if err := db.Create(&user).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}
		} else {
			// Update existing user with Google info if not set
			updated := false
			if user.GoogleID == "" {
				user.GoogleID = userInfo.Id
				updated = true
			}
			if user.Avatar == "" && userInfo.Picture != "" {
				user.Avatar = userInfo.Picture
				updated = true
			}
			if user.FirstName == "" && userInfo.GivenName != "" {
				user.FirstName = userInfo.GivenName
				updated = true
			}
			if user.LastName == "" && userInfo.FamilyName != "" {
				user.LastName = userInfo.FamilyName
				updated = true
			}
			
			if updated {
				if err := db.Save(&user).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
					return
				}
			}
		}

		// Update last login
		now := time.Now()
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

func verifyGoogleIDToken(idToken string) (*oauth2.Userinfo, error) {
	ctx := context.Background()
	
	// Create OAuth2 service
	oauth2Service, err := oauth2.NewService(ctx, option.WithHTTPClient(http.DefaultClient))
	if err != nil {
		return nil, err
	}

	// Verify the token by getting user info
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}

	// Get user info
	userInfoCall := oauth2Service.Userinfo.Get()
	userInfoCall.Context(ctx)
	
	// We need to create a temporary client with the token
	// In a real implementation, you'd use the Google Client Library properly
	userInfo := &oauth2.Userinfo{
		Id:         tokenInfo.UserId,
		Email:      tokenInfo.Email,
		GivenName:  tokenInfo.GivenName,
		FamilyName: tokenInfo.FamilyName,
		Picture:    tokenInfo.Picture,
	}

	return userInfo, nil
}