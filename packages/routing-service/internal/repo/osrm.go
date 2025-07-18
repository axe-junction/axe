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

