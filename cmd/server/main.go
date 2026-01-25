package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Jeya-sGit/CurrencyAggregator/internal/handler"
	"github.com/Jeya-sGit/CurrencyAggregator/internal/providers"
	"github.com/Jeya-sGit/CurrencyAggregator/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	pList := []providers.Provider{
		&providers.MockProvider{},
		&providers.FrankfurterProvider{},
	}
	svc := service.NewAggregatorService(pList)
	h := handler.NewCurrencyHandler(svc)

	r := chi.NewRouter()

	r.Use(middleware.Logger)    // Automatically logs every request to the terminal
	r.Use(middleware.Recoverer) // Prevents the server from crashing if there's a panic
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/rates", h.GetRates)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("🚀 Service started with Chi on http://localhost:8080")
	srv.ListenAndServe()
}
