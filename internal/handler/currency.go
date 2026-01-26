package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Jeya-sGit/CurrencyAggregator/internal/models"
	"github.com/Jeya-sGit/CurrencyAggregator/internal/service"
)

type CurrencyHandler struct {
	service *service.AggregatorService
}

func NewCurrencyHandler(s *service.AggregatorService) *CurrencyHandler {
	return &CurrencyHandler{
		service: s,
	}
}

func sendJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *CurrencyHandler) GetRates(w http.ResponseWriter, r *http.Request) {
	base := r.URL.Query().Get("base")
	target := r.URL.Query().Get("target")
	amountStr := r.URL.Query().Get("amount")

	amount := 1.0
	if amountStr != "" {
		val, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			sendJSONError(w, "invalid amount: must be a numeric value", http.StatusBadRequest)
			return
		}

		if val < 0 {
			sendJSONError(w, "invalid amount: cannot be negative", http.StatusBadRequest)
			return
		}
		amount = val
	}

	req := models.RateRequest{
		BaseCurrency:   base,
		TargetCurrency: target,
		Amount:         amount,
	}

	result, err := h.service.GetAggregateRates(r.Context(), req)
	if err != nil {
		sendJSONError(w, "failed to fetch currency data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {

		return
	}
}
func (h *CurrencyHandler) GetMarketOverview(w http.ResponseWriter, r *http.Request) {
	base := r.URL.Query().Get("base")
	if base == "" {
		base = "INR"
	}

	res, err := h.service.GetMarketData(r.Context(), base)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
