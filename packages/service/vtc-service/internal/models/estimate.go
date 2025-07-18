package models

type GeocodeResponse []struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type EstimateNameRequest struct {
	OriginName      string `json:"originName"`
	DestinationName string `json:"destinationName"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

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

type OSRMResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
	} `json:"routes"`
	Code string `json:"code"`
}
