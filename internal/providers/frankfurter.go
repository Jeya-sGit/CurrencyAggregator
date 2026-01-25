package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Jeya-sGit/CurrencyAggregator/internal/models"
)

type FrankfurterProvider struct{}

type frankfurterResponse struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

func (f *FrankfurterProvider) FetchRate(ctx context.Context, req models.RateRequest) (*models.ProviderResult, error) {
	baseURL := "https://api.frankfurter.app/latest"
	params := url.Values{}
	params.Add("from", req.BaseCurrency)
	params.Add("to", req.TargetCurrency)

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)

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

	var resp frankfurterResponse

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println(resp)

	rate, ok := resp.Rates[req.TargetCurrency]
	if !ok {
		return nil, fmt.Errorf("target currency %s not found in response", req.TargetCurrency)
	}

	return &models.ProviderResult{
		Source: "Frankfurter",
		Rate:   rate,
	}, nil

}
