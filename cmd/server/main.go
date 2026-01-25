package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Jeya-sGit/CurrencyAggregator/internal/models"
	"github.com/Jeya-sGit/CurrencyAggregator/internal/providers"
	"github.com/Jeya-sGit/CurrencyAggregator/internal/service"
)

func main() {

	pList := []providers.Provider{
		&providers.MockProvider{},
	}

	aggregator := service.NewAggregatorService(pList)

	req := models.RateRequest{
		BaseCurrency:   "USD",
		TargetCurrency: "INR",
		Amount:         1.0,
	}

	resp, err := aggregator.GetAggregateRates(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to aggregate rates: %v", err)
	}

	fmt.Println("--- Currency Aggregator Result ---")
	fmt.Printf("Base: %s | Target: %s\n", resp.Base, resp.Target)
	for _, res := range resp.Results {
		fmt.Printf("Source: %-12s | Rate: %.2f\n", res.Source, res.Rate)
	}
	fmt.Printf("Total Latency: %v\n", resp.TotalLatency)
}
