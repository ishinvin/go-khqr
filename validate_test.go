package khqr

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestValidateAccountID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{"valid", "test@bank", nil},
		{"empty", "", ErrAccountIDRequired},
		{"whitespace_only", "   ", ErrAccountIDRequired},
		{"no_at_sign", "noemailformat", ErrAccountIDInvalid},
		{"too_long", strings.Repeat("a", 30) + "@bank", ErrAccountIDTooLong},
		{"max_length", strings.Repeat("a", 27) + "@bank", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateAccountID(tt.id)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateAccountID(%q) = %v, want %v", tt.id, got, tt.wantErr)
			}
		})
	}
}

func TestValidateCurrency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		c       Currency
		wantErr error
	}{
		{"KHR", KHR, nil},
		{"USD", USD, nil},
		{"invalid", Currency(999), ErrInvalidCurrency},
		{"zero", Currency(0), ErrInvalidCurrency},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateCurrency(tt.c)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateCurrency(%d) = %v, want %v", tt.c, got, tt.wantErr)
			}
		})
	}
}

func TestValidateAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		amount   float64
		currency Currency
		wantErr  error
	}{
		{"zero", 0, KHR, nil},
		{"negative", -1, KHR, ErrInvalidAmount},
		{"valid_KHR_integer", 50000, KHR, nil},
		{"KHR_with_decimals", 100.5, KHR, ErrInvalidAmount},
		{"valid_USD", 10.55, USD, nil},
		{"USD_three_decimals", 10.555, USD, ErrInvalidAmount},
		{"USD_one_decimal", 10.5, USD, nil},
		{"USD_integer", 10, USD, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateAmount(tt.amount, tt.currency)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateAmount(%v, %v) = %v, want %v", tt.amount, tt.currency, got, tt.wantErr)
			}
		})
	}
}

func TestValidateTimestamp(t *testing.T) {
	t.Parallel()

	future := time.Now().Add(5 * time.Minute).UnixMilli()
	past := time.Now().Add(-5 * time.Minute).UnixMilli()

	tests := []struct {
		name       string
		expiration int64
		amount     float64
		wantErr    error
	}{
		{"zero_amount_skips", 0, 0, nil},
		{"negative_amount_skips", 0, -1, nil},
		{"missing_with_amount", 0, 100, ErrExpirationRequired},
		{"future", future, 100, nil},
		{"past", past, 100, ErrExpirationInPast},
		{"invalid_length", 123, 100, ErrInvalidTimestamp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateTimestamp(tt.expiration, tt.amount)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateTimestamp(%d, %v) = %v, want %v", tt.expiration, tt.amount, got, tt.wantErr)
			}
		})
	}
}

func TestValidateMerchantName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		val     string
		wantErr error
	}{
		{"valid", "Test Merchant", nil},
		{"empty", "", ErrMerchantNameRequired},
		{"whitespace_only", "   ", ErrMerchantNameRequired},
		{"too_long", strings.Repeat("a", 26), ErrMerchantNameTooLong},
		{"max_length", strings.Repeat("a", 25), nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateMerchantName(tt.val)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateMerchantName(%q) = %v, want %v", tt.val, got, tt.wantErr)
			}
		})
	}
}

func TestValidateMerchantCity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		val     string
		wantErr error
	}{
		{"valid", "Phnom Penh", nil},
		{"empty", "", ErrMerchantCityRequired},
		{"whitespace_only", "   ", ErrMerchantCityRequired},
		{"too_long", strings.Repeat("a", 16), ErrMerchantCityTooLong},
		{"max_length", strings.Repeat("a", 15), nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateMerchantCity(tt.val)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateMerchantCity(%q) = %v, want %v", tt.val, got, tt.wantErr)
			}
		})
	}
}

func TestValidateMerchantCategoryCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		code    string
		wantErr error
	}{
		{"valid_4_digits", "5999", nil},
		{"valid_1_digit", "1", nil},
		{"empty", "", ErrMerchantCategoryCodeRequired},
		{"non_numeric", "abcd", ErrMerchantCategoryCodeInvalid},
		{"too_many_digits", "12345", ErrMerchantCategoryCodeInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateMerchantCategoryCode(tt.code)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateMerchantCategoryCode(%q) = %v, want %v", tt.code, got, tt.wantErr)
			}
		})
	}
}

func TestValidateMerchantID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{"valid", "123456789", nil},
		{"empty", "", ErrMerchantIDRequired},
		{"whitespace_only", "   ", ErrMerchantIDRequired},
		{"too_long", strings.Repeat("a", 33), ErrMerchantIDTooLong},
		{"max_length", strings.Repeat("a", 32), nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateMerchantID(tt.id)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateMerchantID(%q) = %v, want %v", tt.id, got, tt.wantErr)
			}
		})
	}
}

func TestValidateAcquiringBank(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		bank    string
		wantErr error
	}{
		{"valid", "Test Bank", nil},
		{"empty", "", ErrAcquiringBankRequired},
		{"whitespace_only", "   ", ErrAcquiringBankRequired},
		{"too_long", strings.Repeat("a", 33), ErrAcquiringBankTooLong},
		{"max_length", strings.Repeat("a", 32), nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateAcquiringBank(tt.bank)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateAcquiringBank(%q) = %v, want %v", tt.bank, got, tt.wantErr)
			}
		})
	}
}

func TestValidateUPIForGenerate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		upi      string
		currency Currency
		wantErr  error
	}{
		{"empty", "", KHR, nil},
		{"valid_KHR", "upi_account", KHR, nil},
		{"USD_not_supported", "upi_account", USD, ErrUPINotSupportUSD},
		{"too_long", strings.Repeat("a", 100), KHR, ErrUPITooLong},
		{"max_length", strings.Repeat("a", 99), KHR, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateUPIForGenerate(tt.upi, tt.currency)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateUPIForGenerate(%q, %v) = %v, want %v", tt.upi, tt.currency, got, tt.wantErr)
			}
		})
	}
}

func TestValidateLanguageTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		preference string
		nameAlt    string
		cityAlt    string
		wantErr    error
	}{
		{"all_empty", "", "", "", nil},
		{"valid_all_fields", "km", "alt name", "alt city", nil},
		{"valid_without_city_alt", "km", "alt name", "", nil},
		{"name_without_preference", "", "alt name", "", ErrLanguagePreferenceRequired},
		{"city_without_preference", "", "", "alt city", ErrLanguagePreferenceRequired},
		{"preference_without_name", "km", "", "", ErrMerchantNameAltRequired},
		{"preference_too_long", "kmm", "alt name", "", ErrLanguagePreferenceTooLong},
		{"name_alt_too_long", "km", strings.Repeat("a", 26), "", ErrMerchantNameAltTooLong},
		{"city_alt_too_long", "km", "alt name", strings.Repeat("a", 16), ErrMerchantCityAltTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateLanguageTemplate(tt.preference, tt.nameAlt, tt.cityAlt)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateLanguageTemplate(%q, %q, %q) = %v, want %v", tt.preference, tt.nameAlt, tt.cityAlt, got, tt.wantErr)
			}
		})
	}
}

func TestValidateOptionalField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		maxLen  int
		wantErr error
	}{
		{"empty", "", 25, nil},
		{"within_limit", "short", 25, nil},
		{"at_limit", strings.Repeat("a", 25), 25, nil},
		{"exceeds_limit", strings.Repeat("a", 26), 25, ErrBillNumberTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateOptionalField(tt.value, tt.maxLen, ErrBillNumberTooLong)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateOptionalField(%q, %d) = %v, want %v", tt.value, tt.maxLen, got, tt.wantErr)
			}
		})
	}
}

func TestValidateCRC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		crc     string
		wantErr error
	}{
		{"valid", "AB12", nil},
		{"empty", "", ErrCRCRequired},
		{"too_short", "AB", ErrCRCInvalid},
		{"too_long", "AB123", ErrCRCInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateCRC(tt.crc)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateCRC(%q) = %v, want %v", tt.crc, got, tt.wantErr)
			}
		})
	}
}

func TestValidatePayloadFormatIndicator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		pfi     string
		wantErr error
	}{
		{"valid", "01", nil},
		{"empty", "", ErrPayloadFormatIndicatorRequired},
		{"too_long", "012", ErrPayloadFormatIndicatorTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validatePayloadFormatIndicator(tt.pfi)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validatePayloadFormatIndicator(%q) = %v, want %v", tt.pfi, got, tt.wantErr)
			}
		})
	}
}

func TestValidatePointOfInitiationMethod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		method      string
		wantDynamic bool
		wantErr     error
	}{
		{"static", "11", false, nil},
		{"dynamic", "12", true, nil},
		{"too_long", "123", false, ErrPointOfInitiationMethodTooLong},
		{"invalid_value", "99", false, ErrPointOfInitiationMethodInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotDynamic, gotErr := validatePointOfInitiationMethod(tt.method)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("validatePointOfInitiationMethod(%q) error = %v, want %v", tt.method, gotErr, tt.wantErr)
			}
			if gotErr == nil && gotDynamic != tt.wantDynamic {
				t.Errorf("validatePointOfInitiationMethod(%q) dynamic = %v, want %v", tt.method, gotDynamic, tt.wantDynamic)
			}
		})
	}
}

func TestValidateMerchantType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mt      MerchantType
		wantErr error
	}{
		{"individual", Individual, nil},
		{"merchant", Merchant, nil},
		{"empty", MerchantType(""), ErrMerchantTypeRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateMerchantType(tt.mt)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateMerchantType(%q) = %v, want %v", tt.mt, got, tt.wantErr)
			}
		})
	}
}

func TestValidateDecodedMerchantCategoryCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		code    string
		wantErr error
	}{
		{"valid", "5999", nil},
		{"single_digit", "1", nil},
		{"empty", "", ErrMerchantCategoryCodeRequired},
		{"too_long", "12345", ErrMerchantCategoryCodeTooLong},
		{"non_numeric", "abcd", ErrMerchantCategoryCodeInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateDecodedMerchantCategoryCode(tt.code)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateDecodedMerchantCategoryCode(%q) = %v, want %v", tt.code, got, tt.wantErr)
			}
		})
	}
}

func TestValidateTransactionCurrency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tc      string
		wantErr error
	}{
		{"KHR", formatCurrency(KHR), nil},
		{"USD", formatCurrency(USD), nil},
		{"empty", "", ErrCurrencyRequired},
		{"too_long", "1234", ErrTransactionCurrencyTooLong},
		{"too_short", "12", ErrTransactionCurrencyTooLong},
		{"unsupported", "978", ErrInvalidCurrency},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateTransactionCurrency(tt.tc)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateTransactionCurrency(%q) = %v, want %v", tt.tc, got, tt.wantErr)
			}
		})
	}
}

func TestValidateCountryCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cc      string
		wantErr error
	}{
		{"valid", "KH", nil},
		{"empty", "", ErrCountryCodeRequired},
		{"too_long", "KHR", ErrCountryCodeTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateCountryCode(tt.cc)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateCountryCode(%q) = %v, want %v", tt.cc, got, tt.wantErr)
			}
		})
	}
}

func TestValidateUPIForDecode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		upi         string
		txCurrency  string
		countryCode string
		wantErr     error
	}{
		{"empty", "", usdCode, "KH", nil},
		{"valid_KHR", "upi_account", khrCode, "KH", nil},
		{"USD_non_KH_country", "upi_account", usdCode, "US", ErrUPINotSupportUSD},
		{"USD_KH_country", "upi_account", usdCode, "KH", nil},
		{"too_long", strings.Repeat("a", 100), khrCode, "KH", ErrUPITooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := validateUPIForDecode(tt.upi, tt.txCurrency, tt.countryCode)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateUPIForDecode(%q, %q, %q) = %v, want %v", tt.upi, tt.txCurrency, tt.countryCode, got, tt.wantErr)
			}
		})
	}
}

func TestValidateTransactionAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		amount   string
		currency string
		wantErr  error
	}{
		{"empty", "", khrCode, nil},
		{"valid_KHR", "50000", khrCode, nil},
		{"valid_USD", "10.55", usdCode, nil},
		{"KHR_with_decimals", "100.5", khrCode, ErrInvalidAmount},
		{"USD_three_decimals", "10.555", usdCode, ErrInvalidAmount},
		{"negative", "-1", khrCode, ErrInvalidAmount},
		{"non_numeric", "abc", khrCode, ErrInvalidAmount},
		{"zero", "0", khrCode, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := &DecodedData{
				TransactionAmount:   tt.amount,
				TransactionCurrency: tt.currency,
			}
			got := data.validateTransactionAmount()
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateTransactionAmount() amount=%q currency=%q = %v, want %v", tt.amount, tt.currency, got, tt.wantErr)
			}
		})
	}
}

func TestValidateDynamicFields(t *testing.T) {
	t.Parallel()

	futureTS := time.Now().Add(5 * time.Minute).UnixMilli()
	pastTS := time.Now().Add(-5 * time.Minute).UnixMilli()

	tests := []struct {
		name       string
		amount     string
		expiration string
		wantErr    error
	}{
		{"valid", "100", formatTimestamp(futureTS), nil},
		{"empty_amount", "", formatTimestamp(futureTS), ErrInvalidDynamicKHQR},
		{"missing_expiration", "100", "", ErrExpirationRequired},
		{"invalid_expiration_length", "100", "123", ErrInvalidTimestamp},
		{"non_numeric_expiration", "100", "abcdefghijklm", ErrInvalidTimestamp},
		{"expired", "100", formatTimestamp(pastTS), ErrKHQRExpired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := &DecodedData{
				TransactionAmount:   tt.amount,
				ExpirationTimestamp: tt.expiration,
			}
			got := data.validateDynamicFields()
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("validateDynamicFields() amount=%q expiration=%q = %v, want %v", tt.amount, tt.expiration, got, tt.wantErr)
			}
		})
	}
}
