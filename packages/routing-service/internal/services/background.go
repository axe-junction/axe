package services

import (
	"log"
	"time"

	"github.com/axe-junction/axe-server/internal/database"
	"gorm.io/gorm"
)

type BackgroundService struct {
	db              *gorm.DB
	userStopService *UserStopService
	paymentService  *PaymentService
	stopChan        chan bool
}

func NewBackgroundService(
	db *gorm.DB,
	userStopService *UserStopService,
	paymentService *PaymentService,
) *BackgroundService {
	return &BackgroundService{
		db:              db,
		userStopService: userStopService,
		paymentService:  paymentService,
		stopChan:        make(chan bool),
	}
}

// Start begins running background tasks
func (s *BackgroundService) Start() {
	log.Println("Starting background service...")

	go s.runUserStopProcessing()

	go s.runDataCleanup()

	go s.runPaymentMonitoring()
}

func (s *BackgroundService) Stop() {
	log.Println("Stopping background service...")
	close(s.stopChan)
}

func (s *BackgroundService) runUserStopProcessing() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Running user stop demand processing...")
			err := s.userStopService.ProcessHighDemandStops(50.0) 
			if err != nil {
				log.Printf("Error processing high demand stops: %v", err)
			}
		case <-s.stopChan:
			return
		}
	}
}

func (s *BackgroundService) runDataCleanup() {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Running data cleanup...")
			err := database.CleanupExpiredData(s.db)
			if err != nil {
				log.Printf("Error during data cleanup: %v", err)
			}
		case <-s.stopChan:
			return
		}
	}
}

func (s *BackgroundService) runPaymentMonitoring() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// This could monitor payment status, retry failed payments, etc.
			// For now, we'll just log that we're monitoring
			log.Println("Payment monitoring heartbeat...")
		case <-s.stopChan:
			return
		}
	}
}

func (s *BackgroundService) RunOnce() {
	log.Println("Running all background tasks once...")

	if err := s.userStopService.ProcessHighDemandStops(50.0); err != nil {
		log.Printf("Error processing high demand stops: %v", err)
	}

	if err := database.CleanupExpiredData(s.db); err != nil {
		log.Printf("Error during data cleanup: %v", err)
	}

	log.Println("Background tasks completed")
}
