package config

// Config holds application configuration
type Config struct {
	Currency CurrencyConfig
}

// CurrencyConfig holds currency exchange API configuration
type CurrencyConfig struct {
	APIURL string
}

// LoadConfig returns a default configuration
// In production, this would load from environment variables or config files
func LoadConfig() *Config {
	return &Config{
		Currency: CurrencyConfig{
			APIURL: "https://api.exchangerate-api.com/v4/latest", // Default to a free API
		},
	}
}

