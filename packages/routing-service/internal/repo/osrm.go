package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/axe-junction/axe-server/internal/models"
)

type OSRMRepo interface {
	GetRouteBetween(ctx context.Context,fromlon, fromlat, tolon, tolat float64) (*models.OSRMRouteResponse, error)
}
type OsrmReposne struct {
	BaseURL string
	client *http.Client
}

func NewOSRMRepo(baseURL string) *OsrmReposne {

	return &OsrmReposne{
		BaseURL: baseURL,
		client :&http.Client{},
	}
}
func (repo *OsrmReposne) GetRouteBetween(ctx context.Context, fromlon, fromlat, tolon, tolat float64) (*models.OSRMRouteResponse, error) {
	url := repo.BaseURL + "/route/v1/driving/" + 
		fmt.Sprintf("%f,%f;%f,%f", fromlon, fromlat, tolon, tolat) + 
		"?overview=full&geometries=geojson"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get route: %s", resp.Status)
	}

	var routeResponse models.OSRMRouteResponse
	if err := json.NewDecoder(resp.Body).Decode(&routeResponse); err != nil {
		return nil, err
	}

	return &routeResponse, nil
}


