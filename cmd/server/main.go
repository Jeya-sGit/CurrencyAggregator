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
		&providers.ExchangeRate{},
	}
	svc := service.NewAggregatorService(pList)
	h := handler.NewCurrencyHandler(svc)

	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*") // Allows your UI to connect
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	r.Use(middleware.Logger)    // Automatically logs every request to the terminal
	r.Use(middleware.Recoverer) // Prevents the server from crashing if there's a panic
	r.Use(middleware.Timeout(60 * time.Second))

	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Get("/compare", h.GetRates)
	r.Get("/market-overview", h.GetMarketOverview)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("🚀 Service started with Chi on http://localhost:8080")
	srv.ListenAndServe()
}
