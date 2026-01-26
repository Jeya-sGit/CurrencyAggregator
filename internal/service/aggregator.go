package service

import (
	"context"
	"fmt"
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

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resChannel := make(chan *models.ProviderResult, len(s.providers))

	for _, provider := range s.providers {
		go func(p providers.Provider) {
			fmt.Println(p)
			pStart := time.Now()
			res, err := p.FetchRate(ctx, req)
			pDuration := time.Since(pStart).String()

			if err != nil {
				resChannel <- &models.ProviderResult{
					Source:   fmt.Sprintf("%T", p),
					Message:  err.Error(),
					Duration: pDuration,
				}
				return
			}

			res.Duration = pDuration
			resChannel <- res
		}(provider)
	}

	var finalResponse []*models.ProviderResult
	for i := 0; i < len(s.providers); i++ {
		finalResponse = append(finalResponse, <-resChannel)
	}

	return &models.RateResponse{
		Base:         req.BaseCurrency,
		Target:       req.TargetCurrency,
		Results:      finalResponse,
		Timestamp:    time.Now(),
		TotalLatency: time.Since(start).String(),
	}, nil
}

func (s *AggregatorService) GetMarketData(ctx context.Context, base string) (*models.RateResponse, error) {
	for _, p := range s.providers {

		if erProvider, ok := p.(*providers.ExchangeRate); ok {
			req := models.RateRequest{
				BaseCurrency: base,
			}

			res, err := erProvider.FetchMarketRate(ctx, req)
			if err != nil {
				return nil, err
			}

			return &models.RateResponse{
				Base:      base,
				Results:   []*models.ProviderResult{res},
				Timestamp: time.Now(),
			}, nil
		}
	}

	return nil, fmt.Errorf("exchange rate provider not found for market data")
}
