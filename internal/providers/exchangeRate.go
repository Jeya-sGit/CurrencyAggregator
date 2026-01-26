package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Jeya-sGit/CurrencyAggregator/internal/models"
)

type ExchangeRate struct{}

type exchangeRateResponse struct {
	Base  string             `json:"base_code"`
	Date  string             `json:"time_last_update_utc"`
	Rates map[string]float64 `json:"rates"`
}

func helper(ctx context.Context, BaseCurrency string) (*exchangeRateResponse, error) {
	baseURL := "https://open.er-api.com/v6/latest/%s"

	reqURL := fmt.Sprintf(baseURL, BaseCurrency)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpResp, err := http.DefaultClient.Do(httpReq)

	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}

	var resp exchangeRateResponse

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println(resp)

	return &resp, nil

}

func (e *ExchangeRate) FetchRate(ctx context.Context, req models.RateRequest) (*models.ProviderResult, error) {
	RatesData, err := helper(ctx, req.BaseCurrency)
	if err != nil {
		return nil, err
	}
	if RatesData.Rates == nil {
		return nil, fmt.Errorf("rates not found in response")
	}

	return &models.ProviderResult{
		Source: "ExchangeRate-API",
		Rate:   RatesData.Rates[req.TargetCurrency] * req.Amount,
	}, nil
}

func (e *ExchangeRate) FetchMarketRate(ctx context.Context, req models.RateRequest) (*models.ProviderResult, error) {
	RatesData, err := helper(ctx, req.BaseCurrency)
	if err != nil {
		return nil, err
	}
	if RatesData.Rates == nil {
		return nil, fmt.Errorf("rates not found in response")
	}

	return &models.ProviderResult{
		Source:   "ExchangeRate-API",
		AllRates: RatesData.Rates,
	}, nil
}
