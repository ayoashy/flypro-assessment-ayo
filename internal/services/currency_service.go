package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"flypro-assessment-ayo/internal/config"
	"github.com/redis/go-redis/v9"
)

type CurrencyService interface {
	ConvertCurrency(ctx context.Context, amount float64, from, to string) (float64, error)
	GetExchangeRate(ctx context.Context, from, to string) (float64, error)
}

type currencyService struct {
	redisClient *redis.Client
	config      *config.Config
	httpClient  *http.Client
}

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

func NewCurrencyService(redisClient *redis.Client, cfg *config.Config) CurrencyService {
	return &currencyService{
		redisClient: redisClient,
		config:      cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *currencyService) ConvertCurrency(ctx context.Context, amount float64, from, to string) (float64, error) {
	if strings.ToUpper(from) == strings.ToUpper(to) {
		return amount, nil
	}

	rate, err := s.GetExchangeRate(ctx, from, to)
	if err != nil {
		return 0, err
	}

	return amount * rate, nil
}

func (s *currencyService) GetExchangeRate(ctx context.Context, from, to string) (float64, error) {
	// Normalize currency codes
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	if from == to {
		return 1.0, nil
	}

	// Try to get from cache first
	cacheKey := fmt.Sprintf("exchange_rate:%s:%s", from, to)
	cachedRate, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var rate float64
		if _, err := fmt.Sscanf(cachedRate, "%f", &rate); err == nil {
			return rate, nil
		}
	}

	// Fetch from API
	rate, err := s.fetchExchangeRate(ctx, from, to)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch exchange rate: %w", err)
	}

	// Cache for 6 hours
	s.redisClient.Set(ctx, cacheKey, fmt.Sprintf("%f", rate), 6*time.Hour)

	return rate, nil
}

func (s *currencyService) fetchExchangeRate(ctx context.Context, from, to string) (float64, error) {
	// Use ExchangeRate-API (free tier)
	url := fmt.Sprintf("%s/%s", s.config.Currency.APIURL, from)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var data ExchangeRateResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	rate, ok := data.Rates[to]
	if !ok {
		return 0, fmt.Errorf("currency code %s not found in response", to)
	}

	return rate, nil
}