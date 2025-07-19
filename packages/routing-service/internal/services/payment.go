package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/google/uuid"
)

// Enhanced payment service for handling route payments
type PaymentService struct {
	paymentRepo    models.PaymentRepo
	paymentTxRepo  models.PaymentTransactionRepo
	routeCache     models.RouteCache
	routingService *RoutingService
}

func NewPaymentService(
	paymentRepo models.PaymentRepo,
	paymentTxRepo models.PaymentTransactionRepo,
	routeCache models.RouteCache,
	routingService *RoutingService,
) *PaymentService {
	return &PaymentService{
		paymentRepo:    paymentRepo,
		paymentTxRepo:  paymentTxRepo,
		routeCache:     routeCache,
		routingService: routingService,
	}
}

// ProcessRoutePayment handles the complete payment flow for a route
func (s *PaymentService) ProcessRoutePayment(userID uuid.UUID, fromLat, fromLng, toLat, toLng float64, options models.RouteOptions, paymentMethod string) (*models.Payment, error) {
	// First, check if we have this route cached
	routeHash := s.routeCache.GenerateRouteHash(fromLat, fromLng, toLat, toLng, options)
	cachedRoute, err := s.routeCache.GetRoute(routeHash)

	var routeResult *models.RouteResult

	if err == nil {
		// Use cached route for payment
		if err := json.Unmarshal(cachedRoute.RouteData, &routeResult); err != nil {
			// If unmarshal fails, calculate route fresh
			routeResult, err = s.calculateAndCacheRoute(fromLat, fromLng, toLat, toLng, options, routeHash)
			if err != nil {
				return nil, err
			}
		}
	} else {
		// Calculate route fresh and cache it
		routeResult, err = s.calculateAndCacheRoute(fromLat, fromLng, toLat, toLng, options, routeHash)
		if err != nil {
			return nil, err
		}
	}

	return s.processPaymentForRoute(userID, routeResult, paymentMethod)
}

func (s *PaymentService) calculateAndCacheRoute(fromLat, fromLng, toLat, toLng float64, options models.RouteOptions, routeHash string) (*models.RouteResult, error) {
	routeResult, err := s.routingService.GetBestRoute(context.Background(), fromLat, fromLng, toLat, toLng, options)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate route: %w", err)
	}

	// Cache the route for future use
	routeData, _ := json.Marshal(routeResult)
	optionsData, _ := json.Marshal(options)

	cachedRoute := models.CachedRoute{
		RouteHash: routeHash,
		FromLat:   fromLat,
		FromLng:   fromLng,
		ToLat:     toLat,
		ToLng:     toLng,
		RouteData: routeData,
		Options:   optionsData,
	}
	s.routeCache.StoreRoute(cachedRoute)

	return routeResult, nil
}

func (s *PaymentService) processPaymentForRoute(userID uuid.UUID, routeResult *models.RouteResult, paymentMethod string) (*models.Payment, error) {
	// Create payment record
	payment := models.Payment{
		ID:            uuid.New(),
		UserID:        userID,
		RouteID:       uuid.New(), // This could be linked to the cached route ID
		Amount:        routeResult.TotalFare,
		Currency:      routeResult.Currency,
		PaymentMethod: paymentMethod,
		Status:        "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Generate transaction ID for external payment processing
	payment.TransactionID = fmt.Sprintf("axe_%s_%d", payment.ID.String()[:8], time.Now().Unix())

	err := s.paymentRepo.CreatePayment(payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	// Create initial transaction record
	transaction := models.PaymentTransaction{
		PaymentID:       payment.ID,
		TransactionType: "charge",
		Amount:          payment.Amount,
		Currency:        payment.Currency,
		ExternalRef:     payment.TransactionID,
		Status:          "pending",
	}
	s.paymentTxRepo.CreateTransaction(transaction)

	// Process payment asynchronously
	go s.simulatePaymentProcessing(payment.ID, paymentMethod)

	return &payment, nil
}

func (s *PaymentService) simulatePaymentProcessing(paymentID uuid.UUID, paymentMethod string) {
	// Simulate payment processing delay
	time.Sleep(2 * time.Second)

	// Simulate different outcomes based on payment method
	var status string
	var failureReason string

	switch paymentMethod {
	case "cash":
		status = "completed" // Cash payments are always successful
	case "card", "mobile":
		// Simulate 95% success rate
		if time.Now().Unix()%20 == 0 {
			status = "failed"
			failureReason = "insufficient funds"
		} else {
			status = "completed"
		}
	case "subscription":
		// Subscription payments have high success rate
		if time.Now().Unix()%50 == 0 {
			status = "failed"
			failureReason = "subscription expired"
		} else {
			status = "completed"
		}
	default:
		status = "failed"
		failureReason = "unsupported payment method"
	}

	// Update payment status
	s.paymentRepo.UpdatePaymentStatus(paymentID, status)

	// Update transaction status
	transactions, err := s.paymentTxRepo.GetTransactionsByPayment(paymentID)
	if err == nil && len(transactions) > 0 {
		s.paymentTxRepo.UpdateTransactionStatus(transactions[0].ID, status, failureReason)
	}

	log.Printf("Payment %s %s via %s", paymentID.String()[:8], status, paymentMethod)
}

// GetPaymentStatus returns the current status of a payment
func (s *PaymentService) GetPaymentStatus(paymentID uuid.UUID) (string, error) {
	payment, err := s.paymentRepo.GetPaymentByID(paymentID)
	if err != nil {
		return "", err
	}
	return payment.Status, nil
}

// GetUserPaymentHistory returns payment history for a user
func (s *PaymentService) GetUserPaymentHistory(userID uuid.UUID) ([]models.Payment, error) {
	return s.paymentRepo.GetPaymentsByUser(userID)
}

// RefundPayment processes a refund for a payment
func (s *PaymentService) RefundPayment(paymentID uuid.UUID, refundAmount float64, reason string) error {
	payment, err := s.paymentRepo.GetPaymentByID(paymentID)
	if err != nil {
		return err
	}

	if payment.Status != "completed" {
		return fmt.Errorf("cannot refund payment that is not completed")
	}

	if refundAmount > payment.Amount {
		return fmt.Errorf("refund amount cannot exceed original payment amount")
	}

	// Create refund transaction
	transaction := models.PaymentTransaction{
		PaymentID:       paymentID,
		TransactionType: "refund",
		Amount:          refundAmount,
		Currency:        payment.Currency,
		ExternalRef:     fmt.Sprintf("refund_%s_%d", paymentID.String()[:8], time.Now().Unix()),
		Status:          "pending",
	}

	err = s.paymentTxRepo.CreateTransaction(transaction)
	if err != nil {
		return err
	}

	// Process refund asynchronously
	go s.processRefund(transaction.ID, refundAmount == payment.Amount)

	return nil
}

func (s *PaymentService) processRefund(transactionID uuid.UUID, fullRefund bool) {
	// Simulate refund processing
	time.Sleep(1 * time.Second)

	// Refunds are usually successful
	status := "completed"
	if time.Now().Unix()%100 == 0 { // 1% failure rate
		status = "failed"
	}

	s.paymentTxRepo.UpdateTransactionStatus(transactionID, status, "")

	log.Printf("Refund transaction %s %s", transactionID.String()[:8], status)
}
