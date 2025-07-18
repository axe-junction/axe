package models

type OSRMRouteResponse struct {
    Routes []struct {
        Distance  float64 `json:"distance"`
        Duration  float64 `json:"duration"`
        Geometry  struct {
            Coordinates [][]float64 `json:"coordinates"`
            Type        string      `json:"type"`
        } `json:"geometry"`
    } `json:"routes"`
}