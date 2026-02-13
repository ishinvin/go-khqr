package khqr

import (
	"strconv"
	"strings"
)

// formatCurrency returns the ISO 4217 numeric code as a zero-padded 3-digit string.
func formatCurrency(c Currency) string {
	n := int(c)
	if n < 10 { //nolint:mnd // zero-pad single digit
		return "00" + strconv.Itoa(n)
	}
	if n < 100 { //nolint:mnd // zero-pad two digits
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}

// formatAmount formats a transaction amount according to currency rules.
func formatAmount(amount float64, currency Currency) string {
	if currency == KHR {
		return strconv.FormatFloat(amount, 'f', 0, 64)
	}
	// USD: format with 2 decimals, then strip trailing zeros
	s := strconv.FormatFloat(amount, 'f', 2, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

// formatTimestamp formats a Unix millisecond timestamp as a string.
func formatTimestamp(ts int64) string {
	return strconv.FormatInt(ts, 10)
}
