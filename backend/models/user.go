// Database models and API types for PlaySlate
package models

import (
	"time"

	"gorm.io/gorm"
)

// User account with passwordless auth
type User struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Email        string     `json:"email" gorm:"unique;not null"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Avatar       string     `json:"avatar"`
	GoogleID     string     `json:"google_id" gorm:"unique"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	Subscription *Subscription `json:"subscription,omitempty" gorm:"foreignKey:UserID"`
}

// MagicLink for passwordless auth (15min expiry, one-time use)
type MagicLink struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Email     string     `json:"email" gorm:"not null"`
	Token     string     `json:"token" gorm:"unique;not null"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedAt time.Time  `json:"created_at"`
}

// Subscription with Stripe billing ($10/mo or $100/yr)
type Subscription struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	UserID            uint           `json:"user_id" gorm:"not null"`
	StripeCustomerID  string         `json:"stripe_customer_id" gorm:"unique"`
	StripeSubID       string         `json:"stripe_subscription_id" gorm:"unique"`
	PlanType          string         `json:"plan_type"` // "monthly" or "yearly"
	Status            string         `json:"status"`    // "active", "cancelled", etc.
	CurrentPeriodEnd  time.Time      `json:"current_period_end"`
	CancelAtPeriodEnd bool           `json:"cancel_at_period_end"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// API request/response types

type MagicLinkRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type MagicLinkVerifyRequest struct {
	Token string `json:"token" binding:"required"`
}

type GoogleAuthRequest struct {
	IdToken string `json:"id_token" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateSubscriptionRequest struct {
	PlanType        string `json:"plan_type" binding:"required,oneof=monthly yearly"`
	PaymentMethodID string `json:"payment_method_id" binding:"required"`
}