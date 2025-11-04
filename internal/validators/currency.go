package validators

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidCurrencyCodes is a list of valid ISO 4217 currency codes
var ValidCurrencyCodes = map[string]bool{
	"USD": true, "EUR": true, "GBP": true, "JPY": true, "AUD": true,
	"CAD": true, "CHF": true, "CNY": true, "INR": true, "SGD": true,
	"NZD": true, "MXN": true, "HKD": true, "NOK": true, "SEK": true,
	"KRW": true, "TRY": true, "RUB": true, "ZAR": true, "BRL": true,
}

// ValidateCurrency validates if the currency code is valid
func ValidateCurrency(fl validator.FieldLevel) bool {
	currency := strings.ToUpper(fl.Field().String())
	return ValidCurrencyCodes[currency]
}

// RegisterCustomValidators registers all custom validators
func RegisterCustomValidators(v *validator.Validate) error {
	if err := v.RegisterValidation("currency", ValidateCurrency); err != nil {
		return fmt.Errorf("failed to register currency validator: %w", err)
	}
	return nil
}
