package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/paymentmethod"
	"github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/webhook"
	"gorm.io/gorm"

	"playslate-backend/models"
)

var (
	MonthlyPriceID = os.Getenv("STRIPE_MONTHLY_PRICE_ID")
	YearlyPriceID  = os.Getenv("STRIPE_YEARLY_PRICE_ID")
)

func init() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

func HandleCreateSubscription(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var req models.CreateSubscriptionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get user from database
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Check if user already has an active subscription
		var existingSubscription models.Subscription
		if err := db.Where("user_id = ? AND status = ?", userID, "active").First(&existingSubscription).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "User already has an active subscription"})
			return
		}

		// Create or retrieve Stripe customer
		customerID, err := createOrGetStripeCustomer(user.Email, user.FirstName, user.LastName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
			return
		}

		// Attach payment method to customer
		pmParams := &stripe.PaymentMethodAttachParams{
			Customer: stripe.String(customerID),
		}
		if _, err := paymentmethod.Attach(req.PaymentMethodID, pmParams); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to attach payment method"})
			return
		}

		// Set as default payment method
		customerParams := &stripe.CustomerParams{
			InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
				DefaultPaymentMethod: stripe.String(req.PaymentMethodID),
			},
		}
		if _, err := customer.Update(customerID, customerParams); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set default payment method"})
			return
		}

		// Determine price ID based on plan type
		var priceID string
		switch req.PlanType {
		case "monthly":
			priceID = MonthlyPriceID
		case "yearly":
			priceID = YearlyPriceID
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan type"})
			return
		}

		// Create subscription
		subParams := &stripe.SubscriptionParams{
			Customer: stripe.String(customerID),
			Items: []*stripe.SubscriptionItemsParams{
				{
					Price: stripe.String(priceID),
				},
			},
			PaymentBehavior: stripe.String("default_incomplete"),
			PaymentSettings: &stripe.SubscriptionPaymentSettingsParams{
				SaveDefaultPaymentMethod: stripe.String("on_subscription"),
			},
		}

		stripeSub, err := subscription.New(subParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
			return
		}

		// Save subscription to database
		dbSubscription := models.Subscription{
			UserID:           user.ID,
			StripeCustomerID: customerID,
			StripeSubID:      stripeSub.ID,
			PlanType:         req.PlanType,
			Status:           string(stripeSub.Status),
			CurrentPeriodEnd: time.Unix(stripeSub.CurrentPeriodEnd, 0),
		}

		if err := db.Create(&dbSubscription).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save subscription"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"subscription_id": stripeSub.ID,
			"client_secret":   stripeSub.LatestInvoice.PaymentIntent.ClientSecret,
			"status":          stripeSub.Status,
		})
	}
}

func HandleGetSubscriptionStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var subscription models.Subscription
		if err := db.Where("user_id = ?", userID).First(&subscription).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No subscription found"})
			return
		}

		c.JSON(http.StatusOK, subscription)
	}
}

func HandleCancelSubscription(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var subscription models.Subscription
		if err := db.Where("user_id = ? AND status = ?", userID, "active").First(&subscription).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No active subscription found"})
			return
		}

		// Cancel subscription at period end
		params := &stripe.SubscriptionParams{
			CancelAtPeriodEnd: stripe.Bool(true),
		}
		stripeSub, err := subscription.Update(subscription.StripeSubID, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription"})
			return
		}

		// Update database
		subscription.CancelAtPeriodEnd = stripeSub.CancelAtPeriodEnd
		if err := db.Save(&subscription).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":             "Subscription will be cancelled at period end",
			"cancel_at_period_end": stripeSub.CancelAtPeriodEnd,
			"current_period_end":  time.Unix(stripeSub.CurrentPeriodEnd, 0),
		})
	}
}

func HandleStripeWebhook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}

		event, err := webhook.ConstructEvent(body, c.GetHeader("Stripe-Signature"), os.Getenv("STRIPE_WEBHOOK_SECRET"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook signature"})
			return
		}

		switch event.Type {
		case "customer.subscription.updated":
			var subscription stripe.Subscription
			if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse subscription"})
				return
			}
			updateSubscriptionStatus(db, subscription.ID, string(subscription.Status))

		case "customer.subscription.deleted":
			var subscription stripe.Subscription
			if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse subscription"})
				return
			}
			updateSubscriptionStatus(db, subscription.ID, "cancelled")

		case "invoice.payment_succeeded":
			// Handle successful payment
			var invoice stripe.Invoice
			if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse invoice"})
				return
			}
			// Update subscription status if needed
			if invoice.Subscription != nil {
				updateSubscriptionStatus(db, invoice.Subscription.ID, "active")
			}

		case "invoice.payment_failed":
			// Handle failed payment
			var invoice stripe.Invoice
			if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse invoice"})
				return
			}
			if invoice.Subscription != nil {
				updateSubscriptionStatus(db, invoice.Subscription.ID, "past_due")
			}
		}

		c.JSON(http.StatusOK, gin.H{"received": true})
	}
}

func createOrGetStripeCustomer(email, firstName, lastName string) (string, error) {
	// First try to find existing customer
	params := &stripe.CustomerListParams{
		Email: stripe.String(email),
	}
	iter := customer.List(params)
	for iter.Next() {
		return iter.Customer().ID, nil
	}

	// Create new customer
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(fmt.Sprintf("%s %s", firstName, lastName)),
	}
	
	newCustomer, err := customer.New(customerParams)
	if err != nil {
		return "", err
	}

	return newCustomer.ID, nil
}

func updateSubscriptionStatus(db *gorm.DB, stripeSubID, status string) {
	var subscription models.Subscription
	if err := db.Where("stripe_sub_id = ?", stripeSubID).First(&subscription).Error; err != nil {
		return
	}

	subscription.Status = status
	db.Save(&subscription)
}