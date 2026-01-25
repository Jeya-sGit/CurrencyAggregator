package providers

import (
	"context"

	"github.com/Jeya-sGit/CurrencyAggregator/internal/models"
)

type Provider interface {
	FetchRate(context.Context, models.RateRequest) (models.ProviderResult, error)
}
