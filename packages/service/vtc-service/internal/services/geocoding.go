package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/axe-junction/axe-server/vtc-service/internal/config"
)

type GeocodeResponse []struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type GeocodingService struct {
	client *http.Client
}

func NewGeocodingService(cfg *config.Config) *GeocodingService {
	return &GeocodingService{
		client: &http.Client{},
	}
}

func (g *GeocodingService) Geocode(place string) (Coordinates, error) {
	var coords Coordinates
	maxRetries := 3
	baseDelay := 1 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		coords, err := g.geocodeOnce(place)
		if err == nil {
			return coords, nil
		}

		if attempt == maxRetries-1 {
			return coords, err
		}

		delay := baseDelay * time.Duration(1<<uint(attempt))
		time.Sleep(delay)
	}

	return coords, fmt.Errorf("failed to geocode after %d attempts", maxRetries)
}

func (g *GeocodingService) geocodeOnce(place string) (Coordinates, error) {
	var coords Coordinates

	query := url.QueryEscape(place)
	geocodeURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", query)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", geocodeURL, nil)
	if err != nil {
		return coords, fmt.Errorf("failed to create geocoding request: %w", err)
	}

	req.Header.Set("User-Agent", "VTC-Service/1.0 (github.com/axe-junction/axe-server)")

	resp, err := g.client.Do(req)
	if err != nil {
		return coords, fmt.Errorf("failed to make geocoding request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return coords, fmt.Errorf("geocoding service returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return coords, fmt.Errorf("failed to read response body: %w", err)
	}

	var result GeocodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return coords, fmt.Errorf("failed to unmarshal geocoding response: %w", err)
	}

	if len(result) == 0 {
		return coords, fmt.Errorf("place not found")
	}

	lat, err := strconv.ParseFloat(result[0].Lat, 64)
	if err != nil {
		return coords, fmt.Errorf("invalid latitude: %w", err)
	}

	lng, err := strconv.ParseFloat(result[0].Lon, 64)
	if err != nil {
		return coords, fmt.Errorf("invalid longitude: %w", err)
	}

	coords.Lat = lat
	coords.Lng = lng

	return coords, nil
}
