package khqr

import "testing"

func TestFormatCurrency(t *testing.T) {
	tests := []struct {
		currency Currency
		expected string
	}{
		{KHR, "116"},
		{USD, "840"},
		{Currency(1), "001"},
		{Currency(50), "050"},
	}

	for _, tt := range tests {
		result := formatCurrency(tt.currency)
		if result != tt.expected {
			t.Errorf("formatCurrency(%d) = %q, expected %q", int(tt.currency), result, tt.expected)
		}
	}
}

func TestFormatAmount(t *testing.T) {
	tests := []struct {
		amount   float64
		currency Currency
		expected string
	}{
		{50000, KHR, "50000"},
		{100, KHR, "100"},
		{1, USD, "1"},
		{10, USD, "10"},
		{1.12, USD, "1.12"},
		{1.10, USD, "1.1"},
		{100.12, USD, "100.12"},
		{5.0, USD, "5"},
		{10.50, USD, "10.5"},
		{25.50, USD, "25.5"},
		{100.00, USD, "100"},
		{9999999999.99, USD, "9999999999.99"},
	}

	for _, tt := range tests {
		result := formatAmount(tt.amount, tt.currency)
		if result != tt.expected {
			t.Errorf("formatAmount(%v, %v) = %q, expected %q", tt.amount, tt.currency, result, tt.expected)
		}
	}
}

func TestFormatTimestamp(t *testing.T) {
	tests := []struct {
		ts       int64
		expected string
	}{
		{1726821915797, "1726821915797"},
		{0, "0"},
		{1687335902584, "1687335902584"},
	}

	for _, tt := range tests {
		result := formatTimestamp(tt.ts)
		if result != tt.expected {
			t.Errorf("formatTimestamp(%d) = %q, expected %q", tt.ts, result, tt.expected)
		}
	}
}
