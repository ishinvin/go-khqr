package khqr

import "testing"

func TestFormatCurrency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		currency Currency
		want     string
	}{
		{"KHR", KHR, "116"},
		{"USD", USD, "840"},
		{"single_digit", Currency(1), "001"},
		{"two_digits", Currency(50), "050"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := formatCurrency(tt.currency)
			if got != tt.want {
				t.Errorf("formatCurrency(%d) = %q, want %q", int(tt.currency), got, tt.want)
			}
		})
	}
}

func TestFormatAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		amount   float64
		currency Currency
		want     string
	}{
		{"KHR_large", 50000, KHR, "50000"},
		{"KHR_small", 100, KHR, "100"},
		{"USD_integer", 1, USD, "1"},
		{"USD_two_digits", 10, USD, "10"},
		{"USD_two_decimals", 1.12, USD, "1.12"},
		{"USD_trailing_zero", 1.10, USD, "1.1"},
		{"USD_with_decimals", 100.12, USD, "100.12"},
		{"USD_whole_float", 5.0, USD, "5"},
		{"USD_one_decimal", 10.50, USD, "10.5"},
		{"USD_half_cent", 25.50, USD, "25.5"},
		{"USD_zero_decimals", 100.00, USD, "100"},
		{"USD_large", 9999999999.99, USD, "9999999999.99"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := formatAmount(tt.amount, tt.currency)
			if got != tt.want {
				t.Errorf("formatAmount(%v, %v) = %q, want %q", tt.amount, tt.currency, got, tt.want)
			}
		})
	}
}

func TestFormatTimestamp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ts   int64
		want string
	}{
		{"recent", 1726821915797, "1726821915797"},
		{"zero", 0, "0"},
		{"past", 1687335902584, "1687335902584"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := formatTimestamp(tt.ts)
			if got != tt.want {
				t.Errorf("formatTimestamp(%d) = %q, want %q", tt.ts, got, tt.want)
			}
		})
	}
}
