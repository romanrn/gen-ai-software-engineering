package service

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/reznikrn/banking-api/internal/domain"
	portin "github.com/reznikrn/banking-api/internal/ports/in"
)

var accountRegex = regexp.MustCompile(`(?i)^ACC-[A-Z0-9]{5}$`)

var validCurrencies = map[string]struct{}{
	"USD": {}, "EUR": {}, "GBP": {}, "JPY": {}, "CHF": {},
	"AUD": {}, "CAD": {}, "CNY": {}, "HKD": {}, "SGD": {},
	"NOK": {}, "SEK": {}, "DKK": {}, "NZD": {}, "MXN": {},
	"INR": {}, "BRL": {}, "ZAR": {}, "RUB": {}, "TRY": {},
	"PLN": {}, "THB": {}, "IDR": {}, "HUF": {}, "CZK": {},
	"ILS": {}, "CLP": {}, "PHP": {}, "AED": {}, "COP": {},
}

var validTypes = map[string]struct{}{
	"deposit": {}, "withdrawal": {}, "transfer": {},
}

func validateCreateInput(input portin.CreateTransactionInput) error {
	var details []domain.ValidationDetail

	if input.Amount <= 0 {
		details = append(details, domain.ValidationDetail{
			Field:   "amount",
			Message: "must be a positive number",
		})
	} else if !hasTwoDecimalsMax(input.Amount) {
		details = append(details, domain.ValidationDetail{
			Field:   "amount",
			Message: "must have at most 2 decimal places",
		})
	}

	if !accountRegex.MatchString(input.FromAccount) {
		details = append(details, domain.ValidationDetail{
			Field:   "fromAccount",
			Message: "must match format ACC-XXXXX (alphanumeric)",
		})
	}

	txType := strings.ToLower(input.Type)
	if txType == "transfer" && !accountRegex.MatchString(input.ToAccount) {
		details = append(details, domain.ValidationDetail{
			Field:   "toAccount",
			Message: "must match format ACC-XXXXX (alphanumeric) for transfer type",
		})
	}

	if _, ok := validCurrencies[strings.ToUpper(input.Currency)]; !ok {
		details = append(details, domain.ValidationDetail{
			Field:   "currency",
			Message: "invalid ISO 4217 currency code",
		})
	}

	if _, ok := validTypes[txType]; !ok {
		details = append(details, domain.ValidationDetail{
			Field:   "type",
			Message: "must be one of: deposit, withdrawal, transfer",
		})
	}

	if len(details) > 0 {
		return &domain.ErrValidation{Details: details}
	}
	return nil
}

// hasTwoDecimalsMax avoids float64 precision issues by checking the string representation.
func hasTwoDecimalsMax(v float64) bool {
	s := strconv.FormatFloat(v, 'f', -1, 64)
	if dot := strings.IndexByte(s, '.'); dot >= 0 {
		return len(s)-dot-1 <= 2
	}
	return true
}
