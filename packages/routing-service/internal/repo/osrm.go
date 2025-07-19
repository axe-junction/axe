package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/axe-junction/axe-server/internal/models"
)

type OSRMRepo struct {
	baseURL string
}

func NewOSRMRepo(baseURL string) *OSRMRepo {
	return &OSRMRepo{baseURL: baseURL}
}

func (r *OSRMRepo) GetRouteBetween(ctx context.Context, fromLng, fromLat, toLng, toLat float64, profile string) (*models.OSRMResponse, error) {
	url := fmt.Sprintf("%s/route/v1/%s/%.6f,%.6f;%.6f,%.6f?steps=true&geometries=geojson",
		r.baseURL, profile, fromLng, fromLat, toLng, toLat)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var osrmResp models.OSRMResponse
	if err := json.Unmarshal(body, &osrmResp); err != nil {
		return nil, err
	}

	if osrmResp.Code != "Ok" || len(osrmResp.Routes) == 0 {
		return nil, fmt.Errorf("osrm error: %s", osrmResp.Code)
	}

	return &osrmResp, nil
}

// Enhanced OSRM methods for better routing
func (r *OSRMRepo) GetMultiModalRoute(ctx context.Context, fromLng, fromLat, toLng, toLat float64, transportModes []string) (*models.OSRMResponse, error) {
	// For multi-modal routing, we'll use the most appropriate profile
	// This is a simplified implementation - in reality, you'd want to
	// combine different transport modes

	profile := "driving" // Default profile
	if len(transportModes) > 0 {
		// Choose profile based on the first transport mode
		switch transportModes[0] {
		case "walking":
			profile = "walking"
		case "cycling":
			profile = "cycling"
		case "driving":
			profile = "driving"
		default:
			profile = "driving"
		}
	}

	return r.GetRouteBetween(ctx, fromLng, fromLat, toLng, toLat, profile)
}

func (r *OSRMRepo) GetOptimizedRoute(ctx context.Context, coordinates [][]float64, profile string) (*models.OSRMResponse, error) {
	// Build URL for trip optimization (visiting multiple waypoints)
	coordinateStr := ""
	for i, coord := range coordinates {
		if i > 0 {
			coordinateStr += ";"
		}
		coordinateStr += fmt.Sprintf("%.6f,%.6f", coord[0], coord[1])
	}

	url := fmt.Sprintf("%s/trip/v1/%s/%s?steps=true&geometries=geojson",
		r.baseURL, profile, coordinateStr)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var osrmResp models.OSRMResponse
	if err := json.Unmarshal(body, &osrmResp); err != nil {
		return nil, err
	}

	if osrmResp.Code != "Ok" || len(osrmResp.Routes) == 0 {
		return nil, fmt.Errorf("osrm trip optimization error: %s", osrmResp.Code)
	}

	return &osrmResp, nil
}

func (r *OSRMRepo) GetMatrixDurations(ctx context.Context, sources, destinations [][]float64, profile string) ([][]float64, error) {
	// Build coordinate strings
	allCoords := append(sources, destinations...)
	coordinateStr := ""
	for i, coord := range allCoords {
		if i > 0 {
			coordinateStr += ";"
		}
		coordinateStr += fmt.Sprintf("%.6f,%.6f", coord[0], coord[1])
	}

	// Build source and destination indices
	sourceIndices := ""
	for i := 0; i < len(sources); i++ {
		if i > 0 {
			sourceIndices += ";"
		}
		sourceIndices += fmt.Sprintf("%d", i)
	}

	destIndices := ""
	for i := 0; i < len(destinations); i++ {
		if i > 0 {
			destIndices += ";"
		}
		destIndices += fmt.Sprintf("%d", len(sources)+i)
	}

	url := fmt.Sprintf("%s/table/v1/%s/%s?sources=%s&destinations=%s",
		r.baseURL, profile, coordinateStr, sourceIndices, destIndices)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var matrixResp models.OSRMMatrixResponse
	if err := json.Unmarshal(body, &matrixResp); err != nil {
		return nil, err
	}

	return matrixResp.Durations, nil
}
