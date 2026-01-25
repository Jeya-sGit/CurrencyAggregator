package service

import (
	"context"
	"time"

	"github.com/Jeya-sGit/CurrencyAggregator/internal/models"
	"github.com/Jeya-sGit/CurrencyAggregator/internal/providers"
)

type AggregatorService struct {
	providers []providers.Provider
}

func NewAggregatorService(providers []providers.Provider) *AggregatorService {
	return &AggregatorService{
		providers: providers,
	}
}

func (s *AggregatorService) GetAggregateRates(ctx context.Context, req models.RateRequest) (*models.RateResponse, error) {
	start := time.Now()

	resChannel := make(chan models.ProviderResult, len(s.providers))

	for _, provider := range s.providers {
		go func(p providers.Provider) {
			res, _ := p.FetchRate(ctx, req)
			resChannel <- res
		}(provider)
	}

	var finalResponse []models.ProviderResult
	for i := 0; i < len(s.providers); i++ {
		finalResponse = append(finalResponse, <-resChannel)
	}

	return &models.RateResponse{
		Base:         req.BaseCurrency,
		Target:       req.TargetCurrency,
		Results:      finalResponse,
		Timestamp:    time.Now(),
		TotalLatency: time.Since(start),
	}, nil
}
