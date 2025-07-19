package handlers

import (
	"net/http"
	"strconv"

	"github.com/axe-junction/axe-server/gateway/internal/services"
	pb "github.com/axe-junction/axe-server/gateway/pkg/pb"
	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	grpcClients *services.GRPCClients
}

func NewAPIHandler(grpcClients *services.GRPCClients) *APIHandler {
	return &APIHandler{
		grpcClients: grpcClients,
	}
}

// GetBestRoute handles routing requests
func (h *APIHandler) GetBestRoute(c *gin.Context) {
	var req struct {
		FromLat         float64  `json:"from_lat" binding:"required"`
		FromLng         float64  `json:"from_lng" binding:"required"`
		ToLat           float64  `json:"to_lat" binding:"required"`
		ToLng           float64  `json:"to_lng" binding:"required"`
		SafetyFilter    bool     `json:"safety_filter"`
		MinSafetyRating float64  `json:"min_safety_rating"`
		PaymentMethods  []string `json:"payment_methods"`
		MaxWalkingTime  float64  `json:"max_walking_time"`
		IncludePrivate  bool     `json:"include_private"`
		MaxFare         float64  `json:"max_fare"`
		Mode            string   `json:"mode"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"code":    "INVALID_REQUEST",
			"details": err.Error(),
		})
		return
	}

	// Validate coordinates
	if req.FromLat < -90 || req.FromLat > 90 || req.ToLat < -90 || req.ToLat > 90 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid latitude values",
			"code":  "INVALID_COORDINATES",
		})
		return
	}

	if req.FromLng < -180 || req.FromLng > 180 || req.ToLng < -180 || req.ToLng > 180 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid longitude values",
			"code":  "INVALID_COORDINATES",
		})
		return
	}

	if h.grpcClients.RoutingService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Routing service is not available",
			"code":  "SERVICE_UNAVAILABLE",
		})
		return
	}

	// Create gRPC request
	grpcReq := &pb.RouteRequest{
		FromLat: req.FromLat,
		FromLng: req.FromLng,
		ToLat:   req.ToLat,
		ToLng:   req.ToLng,
		Options: &pb.RouteOptions{
			SafetyFilter:    req.SafetyFilter,
			MinSafetyRating: req.MinSafetyRating,
			PaymentMethods:  req.PaymentMethods,
			MaxWalkingTime:  req.MaxWalkingTime,
			IncludePrivate:  req.IncludePrivate,
			MaxFare:         req.MaxFare,
			Mode:            req.Mode,
		},
	}

	// Call routing service
	resp, err := h.grpcClients.RoutingService.GetBestRoute(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get route",
			"code":    "ROUTING_FAILED",
			"details": err.Error(),
		})
		return
	}

	if resp.Error != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": resp.Error,
			"code":  "ROUTING_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"route":   resp.Route,
		"success": true,
	})
}

// GetVTCPrices handles VTC pricing requests
func (h *APIHandler) GetVTCPrices(c *gin.Context) {
	var req struct {
		OriginName      string  `json:"origin_name"`
		DestinationName string  `json:"destination_name"`
		OriginLat       float64 `json:"origin_lat"`
		OriginLng       float64 `json:"origin_lng"`
		DestinationLat  float64 `json:"destination_lat"`
		DestinationLng  float64 `json:"destination_lng"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"code":    "INVALID_REQUEST",
			"details": err.Error(),
		})
		return
	}

	// Validate that we have either names or coordinates
	hasNames := req.OriginName != "" && req.DestinationName != ""
	hasCoords := req.OriginLat != 0 && req.OriginLng != 0 && req.DestinationLat != 0 && req.DestinationLng != 0

	if !hasNames && !hasCoords {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Either location names or coordinates must be provided",
			"code":  "MISSING_LOCATION_DATA",
		})
		return
	}

	if h.grpcClients.VTCService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "VTC service is not available",
			"code":  "SERVICE_UNAVAILABLE",
		})
		return
	}

	// Create gRPC request
	grpcReq := &pb.PriceRequest{
		OriginName:      req.OriginName,
		DestinationName: req.DestinationName,
		OriginLat:       req.OriginLat,
		OriginLng:       req.OriginLng,
		DestinationLat:  req.DestinationLat,
		DestinationLng:  req.DestinationLng,
	}

	// Call VTC service
	resp, err := h.grpcClients.VTCService.GetVTCPrices(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get VTC prices",
			"code":    "VTC_PRICING_FAILED",
			"details": err.Error(),
		})
		return
	}

	if resp.Error != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": resp.Error,
			"code":  "VTC_PRICING_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"estimates": resp.Estimate,
		"success":   true,
	})
}

// GetBestRouteByQuery handles routing requests via query parameters (for GET requests)
func (h *APIHandler) GetBestRouteByQuery(c *gin.Context) {
	fromLatStr := c.Query("from_lat")
	fromLngStr := c.Query("from_lng")
	toLatStr := c.Query("to_lat")
	toLngStr := c.Query("to_lng")

	if fromLatStr == "" || fromLngStr == "" || toLatStr == "" || toLngStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required parameters: from_lat, from_lng, to_lat, to_lng",
			"code":  "MISSING_PARAMETERS",
		})
		return
	}

	fromLat, err := strconv.ParseFloat(fromLatStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid from_lat parameter",
			"code":  "INVALID_PARAMETER",
		})
		return
	}

	fromLng, err := strconv.ParseFloat(fromLngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid from_lng parameter",
			"code":  "INVALID_PARAMETER",
		})
		return
	}

	toLat, err := strconv.ParseFloat(toLatStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid to_lat parameter",
			"code":  "INVALID_PARAMETER",
		})
		return
	}

	toLng, err := strconv.ParseFloat(toLngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid to_lng parameter",
			"code":  "INVALID_PARAMETER",
		})
		return
	}

	// Parse optional parameters
	safetyFilter := c.Query("safety_filter") == "true"
	minSafetyRating, _ := strconv.ParseFloat(c.DefaultQuery("min_safety_rating", "3.0"), 64)
	maxWalkingTime, _ := strconv.ParseFloat(c.DefaultQuery("max_walking_time", "900"), 64)
	includePrivate := c.Query("include_private") == "true"
	maxFare, _ := strconv.ParseFloat(c.DefaultQuery("max_fare", "0"), 64)
	mode := c.DefaultQuery("mode", "balanced")

	if h.grpcClients.RoutingService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Routing service is not available",
			"code":  "SERVICE_UNAVAILABLE",
		})
		return
	}

	// Create gRPC request
	grpcReq := &pb.RouteRequest{
		FromLat: fromLat,
		FromLng: fromLng,
		ToLat:   toLat,
		ToLng:   toLng,
		Options: &pb.RouteOptions{
			SafetyFilter:    safetyFilter,
			MinSafetyRating: minSafetyRating,
			MaxWalkingTime:  maxWalkingTime,
			IncludePrivate:  includePrivate,
			MaxFare:         maxFare,
			Mode:            mode,
		},
	}

	// Call routing service
	resp, err := h.grpcClients.RoutingService.GetBestRoute(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get route",
			"code":    "ROUTING_FAILED",
			"details": err.Error(),
		})
		return
	}

	if resp.Error != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": resp.Error,
			"code":  "ROUTING_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"route":   resp.Route,
		"success": true,
	})
}

// HealthCheck returns the health status of the gateway
func (h *APIHandler) HealthCheck(c *gin.Context) {
	status := gin.H{
		"status": "healthy",
		"services": gin.H{
			"routing": h.grpcClients.RoutingService != nil,
			"vtc":     h.grpcClients.VTCService != nil,
		},
	}

	c.JSON(http.StatusOK, status)
}
