package providers

import (
	"context"

	"github.com/Jeya-sGit/CurrencyAggregator/internal/models"
)

type MockProvider struct{}

func (m *MockProvider) FetchRate(ctx context.Context, req models.RateRequest) (models.ProviderResult, error) {

	return models.ProviderResult{
		Source: "MockProvider",
		Rate:   84.2,
	}, nil
}
