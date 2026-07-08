package models

type ErrorResponse struct {
	Error string `json:"error"`
}

type AnalyticsSummary struct {
	AvailabilityPercent float64 `json:"availability_percent"`
	SLAPercent          float64 `json:"sla_percent"`
}
