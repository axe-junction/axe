package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/axe-junction/axe-server/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetRouteRequest struct {
	FromLat float64 `json:"from_lat" binding:"required"`
	FromLng float64 `json:"from_lng" binding:"required"`
	ToLat   float64 `json:"to_lat" binding:"required"`
	ToLng   float64 `json:"to_lng" binding:"required"`
	// Enhanced route options
	SafetyFilter    bool     `json:"safety_filter"`
	MinSafetyRating float64  `json:"min_safety_rating"`
	PaymentMethods  []string `json:"payment_methods"`
	MaxWalkingTime  float64  `json:"max_walking_time"`
	IncludePrivate  bool     `json:"include_private"`
	MaxFare         float64  `json:"max_fare"`
}

// Enhanced handler with safety filtering and payment options
func GetRouteHandler(routeService services.RoutingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req GetRouteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set default values
		options := models.RouteOptions{
			SafetyFilter:    req.SafetyFilter,
			MinSafetyRating: req.MinSafetyRating,
			PaymentMethods:  req.PaymentMethods,
			MaxWalkingTime:  req.MaxWalkingTime,
			IncludePrivate:  req.IncludePrivate,
			MaxFare:         req.MaxFare,
		}

		// Set defaults if not provided
		if options.MinSafetyRating == 0 {
			options.MinSafetyRating = 3.0 // Default minimum safety rating
		}
		if options.MaxWalkingTime == 0 {
			options.MaxWalkingTime = 900 // 15 minutes default
		}

		resp, err := routeService.GetBestRoute(
			c.Request.Context(),
			req.FromLat, req.FromLng,
			req.ToLat, req.ToLng,
			options,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

// Handler for multiple route options
func GetMultipleRoutesHandler(routeService services.RoutingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req GetRouteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		options := models.RouteOptions{
			SafetyFilter:    req.SafetyFilter,
			MinSafetyRating: req.MinSafetyRating,
			PaymentMethods:  req.PaymentMethods,
			MaxWalkingTime:  req.MaxWalkingTime,
			IncludePrivate:  req.IncludePrivate,
			MaxFare:         req.MaxFare,
		}

		if options.MinSafetyRating == 0 {
			options.MinSafetyRating = 3.0
		}

		routes, err := routeService.GetMultipleRoutes(
			c.Request.Context(),
			req.FromLat, req.FromLng,
			req.ToLat, req.ToLng,
			options,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"routes": routes})
	}
}

// Handler for adding user-contributed bus stops
type AddUserStopRequest struct {
	Name           string  `json:"name" binding:"required"`
	Latitude       float64 `json:"latitude" binding:"required"`
	Longitude      float64 `json:"longitude" binding:"required"`
	Description    string  `json:"description"`
	NearbyLandmark string  `json:"nearby_landmark"`
}

func AddUserContributedStopHandler(routeService services.RoutingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AddUserStopRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get user ID from context (assuming auth middleware sets this)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		stop := models.UserContributedStop{
			ID:             uuid.New(),
			Name:           req.Name,
			Latitude:       req.Latitude,
			Longitude:      req.Longitude,
			ContributorID:  userID.(uuid.UUID),
			Description:    req.Description,
			NearbyLandmark: req.NearbyLandmark,
			Status:         "pending",
			RequestCount:   1,
			DemandScore:    1.0,
		}

		err := routeService.AddUserContributedStop(stop)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User-contributed stop added successfully",
			"stop":    stop,
		})
	}
}

// Handler for requesting demand for a user-contributed stop
type RequestStopRequest struct {
	StopID  string  `json:"stop_id" binding:"required"`
	FromLat float64 `json:"from_lat" binding:"required"`
	FromLng float64 `json:"from_lng" binding:"required"`
	ToLat   float64 `json:"to_lat" binding:"required"`
	ToLng   float64 `json:"to_lng" binding:"required"`
}

func RequestUserStopHandler(routeService services.RoutingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RequestStopRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		stopID, err := uuid.Parse(req.StopID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stop ID"})
			return
		}

		err = routeService.RequestUserStop(
			stopID,
			userID.(uuid.UUID),
			req.FromLat, req.FromLng,
			req.ToLat, req.ToLng,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Stop demand request recorded"})
	}
}

// Handler for getting nearby user-contributed stops
func GetNearbyUserStopsHandler(routeService services.RoutingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		lat := c.Query("lat")
		lng := c.Query("lng")
		radius := c.DefaultQuery("radius", "1000")

		if lat == "" || lng == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lng parameters are required"})
			return
		}

		// Parse parameters
		latVal, err := parseFloat(lat)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
			return
		}
		lngVal, err := parseFloat(lng)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
			return
		}
		radiusVal, err := parseInt(radius)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid radius"})
			return
		}

		stops, err := routeService.GetNearbyUserStops(latVal, lngVal, radiusVal)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"stops": stops})
	}
}

// Handler for processing route payments
type ProcessPaymentRequest struct {
	RouteID       string `json:"route_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"`
}

func ProcessRoutePaymentHandler(routeService services.RoutingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ProcessPaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// This would need to be implemented with route caching or route reconstruction
		// For now, we'll return a placeholder response
		c.JSON(http.StatusOK, gin.H{
			"message":        "Payment processing endpoint - implementation needed",
			"user_id":        userID,
			"route_id":       req.RouteID,
			"payment_method": req.PaymentMethod,
		})
	}
}

// Enhanced payment processing handler with proper route caching
func ProcessRoutePaymentEnhancedHandler(paymentService *services.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ProcessPaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		// Get route parameters from request
		fromLat := parseFloatFromQuery(c, "from_lat")
		fromLng := parseFloatFromQuery(c, "from_lng")
		toLat := parseFloatFromQuery(c, "to_lat")
		toLng := parseFloatFromQuery(c, "to_lng")

		if fromLat == 0 || fromLng == 0 || toLat == 0 || toLng == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route coordinates"})
			return
		}

		// Default route options
		options := models.RouteOptions{
			SafetyFilter:    true,
			MinSafetyRating: 3.0,
			MaxWalkingTime:  900, // 15 minutes
		}

		payment, err := paymentService.ProcessRoutePayment(
			userID.(uuid.UUID),
			fromLat, fromLng, toLat, toLng,
			options,
			req.PaymentMethod,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"payment": payment,
			"message": "Payment processing initiated",
		})
	}
}

// Add review handler
type AddReviewRequest struct {
	StationID    *uuid.UUID `json:"station_id,omitempty"`
	LineID       *uuid.UUID `json:"line_id,omitempty"`
	Rating       float64    `json:"rating" binding:"required,min=1,max=5"`
	SafetyRating float64    `json:"safety_rating" binding:"required,min=1,max=5"`
	Comment      string     `json:"comment"`
	ReviewType   string     `json:"review_type" binding:"required"`
}

func AddReviewHandler(reviewService *services.ReviewService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AddReviewRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		review := models.Review{
			StationID:    req.StationID,
			LineID:       req.LineID,
			Rating:       req.Rating,
			SafetyRating: req.SafetyRating,
			Comment:      req.Comment,
			ReviewType:   req.ReviewType,
		}

		err := reviewService.AddReview(userID.(uuid.UUID), review)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Review added successfully",
			"review":  review,
		})
	}
}

// Get reviews for station
func GetStationReviewsHandler(reviewService *services.ReviewService) gin.HandlerFunc {
	return func(c *gin.Context) {
		stationIDStr := c.Param("stationId")
		stationID, err := uuid.Parse(stationIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid station ID"})
			return
		}

		reviews, err := reviewService.GetStationReviews(stationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rating, safety, count, err := reviewService.GetStationRating(stationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"reviews":        reviews,
			"average_rating": rating,
			"safety_rating":  safety,
			"review_count":   count,
		})
	}
}

// Get payment status
func GetPaymentStatusHandler(paymentService *services.PaymentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		paymentIDStr := c.Param("paymentId")
		paymentID, err := uuid.Parse(paymentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
			return
		}

		status, err := paymentService.GetPaymentStatus(paymentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"payment_id": paymentID,
			"status":     status,
		})
	}
}

// Enhanced user stop management
func ManageUserStopsHandler(userStopService *services.UserStopService) gin.HandlerFunc {
	return func(c *gin.Context) {
		action := c.Query("action")

		switch action {
		case "process_demand":
			err := userStopService.ProcessHighDemandStops(50.0)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "High demand stops processed"})

		case "statistics":
			stats, err := userStopService.GetStopStatistics()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, stats)

		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		}
	}
}

func GetFareEstimateHandler(routeService services.RoutingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		transportType := c.Query("transport_type")
		distanceStr := c.Query("distance")

		if transportType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "transport_type is required"})
			return
		}

		distance := 5.0 // Default 5km
		if distanceStr != "" {
			if d, err := strconv.ParseFloat(distanceStr, 64); err == nil {
				distance = d
			}
		}

		// Get fare estimate
		fareInfo := routeService.GetFareEstimate(transportType, distance)

		c.JSON(http.StatusOK, gin.H{
			"transport_type": fareInfo.TransportType,
			"distance":       fareInfo.Distance,
			"base_fare":      fareInfo.BaseFare,
			"distance_fare":  fareInfo.DistanceFare,
			"total_fare":     fareInfo.TotalFare,
			"currency":       fareInfo.Currency,
			"estimated_cost": fmt.Sprintf("%.0f DZD", fareInfo.TotalFare),
		})
	}
}

func GetRouteWithModeHandler(routeService services.RoutingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req GetRouteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		mode := c.Query("mode") // economy, fastest, safest, balanced, comfort
		if mode == "" {
			mode = "balanced"
		}

		// Set default values
		options := models.RouteOptions{
			Mode:            mode,
			SafetyFilter:    req.SafetyFilter,
			MinSafetyRating: req.MinSafetyRating,
			PaymentMethods:  req.PaymentMethods,
			MaxWalkingTime:  req.MaxWalkingTime,
			IncludePrivate:  req.IncludePrivate,
			MaxFare:         req.MaxFare,
		}

		if options.MinSafetyRating == 0 {
			options.MinSafetyRating = 3.0 // Default minimum safety rating
		}
		if options.MaxWalkingTime == 0 {
			options.MaxWalkingTime = 900 // 15 minutes default
		}

		route, err := routeService.GetBestRoute(c.Request.Context(), req.FromLat, req.FromLng, req.ToLat, req.ToLng, options)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := gin.H{
			"route": route,
			"mode":  mode,
			"optimization_info": gin.H{
				"optimized_for": mode,
				"total_fare":    route.TotalFare,
				"total_time":    route.TotalTime,
				"safety_rating": route.OverallSafetyRating,
				"walking_time":  route.Summary.WalkingDuration,
				"transfers":     route.Summary.TransfersCount,
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func parseFloatFromQuery(c *gin.Context, key string) float64 {
	val := c.Query(key)
	if val == "" {
		return 0
	}
	result, _ := parseFloat(val)
	return result
}
