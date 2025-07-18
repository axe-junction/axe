package services

import (
	"fmt"
	"math"

	"github.com/axe-junction/axe-server/vtc-service/internal/config"
)

type PriceEstimate struct {
	Price int    `json:"price"`
	ETA   string `json:"eta"`
}

type RangeEstimate struct {
	PriceRange string `json:"priceRange"`
	ETA        string `json:"eta"`
}

type EstimateResponse struct {
	ServiceA PriceEstimate `json:"serviceA"`
	ServiceB RangeEstimate `json:"serviceB"`
	ServiceC PriceEstimate `json:"serviceC"`
	ServiceD PriceEstimate `json:"serviceD"`
}

type PricingService struct {
	rates config.RatesConfig
}

func NewPricingService(cfg *config.Config) *PricingService {
	return &PricingService{
		rates: cfg.Rates,
	}
}

func (p *PricingService) CalculateEstimates(distance, duration float64) EstimateResponse {
	distanceKm := distance / 1000.0
	durationMin := int(duration / 60)
	eta := fmt.Sprintf("%d min", durationMin)

	return EstimateResponse{
		ServiceA: PriceEstimate{
			Price: int(math.Round(distanceKm * p.rates.ServiceA)),
			ETA:   eta,
		},
		ServiceB: RangeEstimate{
			PriceRange: fmt.Sprintf("%d–%d",
				int(distanceKm*p.rates.ServiceBMin),
				int(distanceKm*p.rates.ServiceBMax)),
			ETA: eta,
		},
		ServiceC: PriceEstimate{
			Price: int(math.Round(distanceKm * p.rates.ServiceC)),
			ETA:   eta,
		},
		ServiceD: PriceEstimate{
			Price: int(math.Round(distanceKm * p.rates.ServiceD)),
			ETA:   eta,
		},
	}
}
