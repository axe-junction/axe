package repo

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) CreatePayment(payment models.Payment) error {
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()
	return r.db.Create(&payment).Error
}

func (r *PaymentRepo) GetPaymentsByUser(userID uuid.UUID) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&payments).Error
	return payments, err
}

func (r *PaymentRepo) GetPaymentByID(paymentID uuid.UUID) (models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("id = ?", paymentID).First(&payment).Error
	return payment, err
}

func (r *PaymentRepo) UpdatePaymentStatus(paymentID uuid.UUID, status string) error {
	return r.db.Model(&models.Payment{}).
		Where("id = ?", paymentID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *PaymentRepo) CalculateRouteFare(segments []models.RouteSegment) (float64, error) {
	totalFare := 0.0

	for _, segment := range segments {
		if segment.PublicRoute != nil {
			// Calculate fare based on line pricing
			baseFare := segment.PublicRoute.BaseFare
			distanceFare := segment.PublicRoute.FarePerKM * (segment.Distance / 1000) // Convert meters to km
			segmentFare := baseFare + distanceFare

			// Apply any discounts or special pricing rules here
			totalFare += segmentFare
		}
		// Walking segments are free
	}

	// Apply any route-level discounts or surcharges
	return totalFare, nil
}

// PaymentTransactionRepo implementation
type PaymentTransactionRepo struct {
	db *gorm.DB
}

func NewPaymentTransactionRepo(db *gorm.DB) *PaymentTransactionRepo {
	return &PaymentTransactionRepo{db: db}
}

func (r *PaymentTransactionRepo) CreateTransaction(transaction models.PaymentTransaction) error {
	transaction.ID = uuid.New()
	transaction.CreatedAt = time.Now()
	return r.db.Create(&transaction).Error
}

func (r *PaymentTransactionRepo) GetTransactionsByPayment(paymentID uuid.UUID) ([]models.PaymentTransaction, error) {
	var transactions []models.PaymentTransaction
	err := r.db.Where("payment_id = ?", paymentID).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *PaymentTransactionRepo) UpdateTransactionStatus(transactionID uuid.UUID, status, failureReason string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == "completed" {
		now := time.Now()
		updates["processed_at"] = &now
	}

	if failureReason != "" {
		updates["failure_reason"] = failureReason
	}

	return r.db.Model(&models.PaymentTransaction{}).
		Where("id = ?", transactionID).
		Updates(updates).Error
}

func (r *PaymentTransactionRepo) GetPendingTransactions() ([]models.PaymentTransaction, error) {
	var transactions []models.PaymentTransaction
	err := r.db.Where("status = ?", "pending").Find(&transactions).Error
	return transactions, err
}

// RouteCache implementation
type RouteCacheRepo struct {
	db *gorm.DB
}

func NewRouteCacheRepo(db *gorm.DB) *RouteCacheRepo {
	return &RouteCacheRepo{db: db}
}

func (r *RouteCacheRepo) StoreRoute(route models.CachedRoute) error {
	route.ID = uuid.New()
	route.CreatedAt = time.Now()
	route.ExpiresAt = time.Now().Add(24 * time.Hour) // Cache for 24 hours
	return r.db.Create(&route).Error
}

func (r *RouteCacheRepo) GetRoute(routeHash string) (*models.CachedRoute, error) {
	var route models.CachedRoute
	err := r.db.Where("route_hash = ? AND expires_at > ?", routeHash, time.Now()).First(&route).Error
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *RouteCacheRepo) DeleteExpiredRoutes() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.CachedRoute{}).Error
}

func (r *RouteCacheRepo) GenerateRouteHash(fromLat, fromLng, toLat, toLng float64, options models.RouteOptions) string {
	data := struct {
		FromLat float64
		FromLng float64
		ToLat   float64
		ToLng   float64
		Options models.RouteOptions
	}{
		FromLat: fromLat,
		FromLng: fromLng,
		ToLat:   toLat,
		ToLng:   toLng,
		Options: options,
	}

	jsonData, _ := json.Marshal(data)
	hash := md5.Sum(jsonData)
	return fmt.Sprintf("%x", hash)
}
