package models

import "time"

type RateRequest struct {
	BaseCurrency   string  `json:"base_currency"`
	TargetCurrency string  `json:"target_currency"`
	Amount         float64 `json:"amount"`
}

type ProviderResult struct {
	Source   string  `json:"source"`
	Rate     float64 `json:"rate"`
	Duration string  `json:"duration"`
	Message  string  `json:"error,omitempty"`
}

type RateResponse struct {
	Base         string            `json:"base"`
	Target       string            `json:"target"`
	Results      []*ProviderResult `json:"results"`
	Timestamp    time.Time         `json:"timestamp"`
	TotalLatency string            `json:"total_latency"`
}
