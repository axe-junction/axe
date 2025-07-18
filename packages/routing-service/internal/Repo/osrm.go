package repo

import (
	"context"
	"net/http"
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
	
}


