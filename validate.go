package khqr

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var accountIDRegex = regexp.MustCompile(`^[^@]+@[^@]+$`)
var merchantCategoryCodeRegex = regexp.MustCompile(`^\d{1,4}$`)

var (
	khrCode = formatCurrency(KHR)
	usdCode = formatCurrency(USD)
)

type checker struct {
	err error
}

func (c *checker) check(fn func() error) {
	if c.err != nil {
		return
	}
	c.err = fn()
}

func (c *checker) optional(value string, maxLen int, errTooLong *Error) {
	c.check(func() error { return validateOptionalField(value, maxLen, errTooLong) })
}

func validateOptionalField(value string, maxLen int, errTooLong *Error) error {
	if value != "" && utf8.RuneCountInString(value) > maxLen {
		return errTooLong
	}
	return nil
}

// --- Field validators (generate) ---

func validateUPIForGenerate(upi string, currency Currency) error {
	if strings.TrimSpace(upi) != "" && currency == USD {
		return ErrUPINotSupportUSD
	}
	if strings.TrimSpace(upi) != "" && utf8.RuneCountInString(upi) > maxUPILength {
		return ErrUPITooLong
	}
	return nil
}

func validateAccountID(id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrAccountIDRequired
	}
	if utf8.RuneCountInString(id) > maxAccountIDLength {
		return ErrAccountIDTooLong
	}
	if !accountIDRegex.MatchString(id) {
		return ErrAccountIDInvalid
	}
	return nil
}

func validateCurrency(c Currency) error {
	if c != KHR && c != USD {
		return ErrInvalidCurrency
	}
	return nil
}

func validateAmount(amount float64, currency Currency) error {
	if amount < 0 {
		return ErrInvalidAmount
	}
	if amount == 0 {
		return nil
	}

	switch currency {
	case KHR:
		if math.Abs(amount-math.Floor(amount)) > 1e-9 { //nolint:mnd // float epsilon
			return ErrInvalidAmount
		}
	case USD:
		rounded := math.Round(amount*100) / 100 //nolint:mnd // cents multiplier
		if math.Abs(amount-rounded) > 1e-9 {    //nolint:mnd // float epsilon
			return ErrInvalidAmount
		}
	}

	amountStr := formatAmount(amount, currency)
	if len(amountStr) > maxAmountLength {
		return ErrInvalidAmount
	}

	return nil
}

func validateTimestamp(expiration int64, amount float64) error {
	if amount <= 0 {
		return nil
	}

	if expiration == 0 {
		return ErrExpirationRequired
	}

	expStr := formatTimestamp(expiration)
	if len(expStr) != 13 { //nolint:mnd // timestamp string length
		return ErrInvalidTimestamp
	}

	now := time.Now().UnixMilli()
	if expiration <= now {
		return ErrExpirationInPast
	}

	return nil
}

func validateMerchantName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrMerchantNameRequired
	}
	if utf8.RuneCountInString(name) > maxMerchantNameLength {
		return ErrMerchantNameTooLong
	}
	return nil
}

func validateMerchantCity(city string) error {
	if strings.TrimSpace(city) == "" {
		return ErrMerchantCityRequired
	}
	if utf8.RuneCountInString(city) > maxMerchantCityLength {
		return ErrMerchantCityTooLong
	}
	return nil
}

func validateMerchantCategoryCode(code string) error {
	if strings.TrimSpace(code) == "" {
		return ErrMerchantCategoryCodeRequired
	}
	if !merchantCategoryCodeRegex.MatchString(code) {
		return ErrMerchantCategoryCodeInvalid
	}
	return nil
}

func validateMerchantID(id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrMerchantIDRequired
	}
	if utf8.RuneCountInString(id) > maxMerchantIDLength {
		return ErrMerchantIDTooLong
	}
	return nil
}

func validateAcquiringBank(bank string) error {
	if strings.TrimSpace(bank) == "" {
		return ErrAcquiringBankRequired
	}
	if utf8.RuneCountInString(bank) > maxAcquiringBankLength {
		return ErrAcquiringBankTooLong
	}
	return nil
}

func validateLanguageTemplate(preference, nameAlt, cityAlt string) error {
	hasAlt := strings.TrimSpace(nameAlt) != "" || strings.TrimSpace(cityAlt) != ""
	if hasAlt && strings.TrimSpace(preference) == "" {
		return ErrLanguagePreferenceRequired
	}
	if strings.TrimSpace(preference) != "" && utf8.RuneCountInString(preference) != languagePreferenceLength {
		return ErrLanguagePreferenceTooLong
	}
	if strings.TrimSpace(preference) != "" && strings.TrimSpace(nameAlt) == "" {
		return ErrMerchantNameAltRequired
	}
	if strings.TrimSpace(nameAlt) != "" && utf8.RuneCountInString(strings.TrimSpace(nameAlt)) > maxMerchantNameAltLength {
		return ErrMerchantNameAltTooLong
	}
	if strings.TrimSpace(cityAlt) != "" && utf8.RuneCountInString(strings.TrimSpace(cityAlt)) > maxMerchantCityAltLength {
		return ErrMerchantCityAltTooLong
	}
	return nil
}

// --- Field validators (decode) ---

func validateCRC(crc string) error {
	if strings.TrimSpace(crc) == "" {
		return ErrCRCRequired
	}
	if len(crc) != 4 { //nolint:mnd // CRC hex length
		return ErrCRCInvalid
	}
	return nil
}

func validatePayloadFormatIndicator(pfi string) error {
	if strings.TrimSpace(pfi) == "" {
		return ErrPayloadFormatIndicatorRequired
	}
	if len(pfi) != 2 {
		return ErrPayloadFormatIndicatorTooLong
	}
	return nil
}

func validatePointOfInitiationMethod(method string) (bool, error) {
	if strings.TrimSpace(method) != "" && len(method) > 2 {
		return false, ErrPointOfInitiationMethodTooLong
	}
	if method != staticQR && method != dynamicQR {
		return false, ErrPointOfInitiationMethodInvalid
	}
	return method == dynamicQR, nil
}

func validateMerchantType(mt MerchantType) error {
	if strings.TrimSpace(string(mt)) == "" {
		return ErrMerchantTypeRequired
	}
	return nil
}

func validateDecodedMerchantCategoryCode(code string) error {
	if strings.TrimSpace(code) == "" {
		return ErrMerchantCategoryCodeRequired
	}
	if len(code) > 4 { //nolint:mnd // MCC is up to 4 digits
		return ErrMerchantCategoryCodeTooLong
	}
	if !merchantCategoryCodeRegex.MatchString(code) {
		return ErrMerchantCategoryCodeInvalid
	}
	return nil
}

func validateTransactionCurrency(tc string) error {
	if strings.TrimSpace(tc) == "" {
		return ErrCurrencyRequired
	}
	if len(tc) != 3 {
		return ErrTransactionCurrencyTooLong
	}
	if tc != khrCode && tc != usdCode {
		return ErrInvalidCurrency
	}
	return nil
}

func validateCountryCode(cc string) error {
	if strings.TrimSpace(cc) == "" {
		return ErrCountryCodeRequired
	}
	if len(cc) > 2 {
		return ErrCountryCodeTooLong
	}
	return nil
}

func validateUPIForDecode(upi, transactionCurrency, countryCode string) error {
	if strings.TrimSpace(upi) == "" {
		return nil
	}
	if transactionCurrency == usdCode && countryCode != defaultCountryCode {
		return ErrUPINotSupportUSD
	}
	if utf8.RuneCountInString(upi) > maxUPILength {
		return ErrUPITooLong
	}
	return nil
}

// --- Validator ---

func (info *IndividualInfo) validate() error {
	var c checker
	c.check(func() error { return validateUPIForGenerate(info.UPIAccountInfo, info.Currency) })
	c.check(func() error { return validateAccountID(info.BakongAccountID) })
	c.check(func() error { return validateCurrency(info.Currency) })
	c.check(func() error { return validateAmount(info.Amount, info.Currency) })
	c.check(func() error { return validateMerchantName(info.MerchantName) })
	c.check(func() error { return validateMerchantCity(info.MerchantCity) })
	c.optional(info.AccountInfo, maxAccountIDLength, ErrAccountInfoTooLong)
	c.optional(info.AcquiringBank, maxAcquiringBankLength, ErrAcquiringBankTooLong)
	c.check(func() error { return validateMerchantCategoryCode(info.MerchantCategoryCode) })
	c.optional(info.TerminalLabel, maxTerminalLabelLength, ErrTerminalLabelTooLong)
	c.optional(info.StoreLabel, maxStoreLabelLength, ErrStoreLabelTooLong)
	c.optional(info.BillNumber, maxBillNumberLength, ErrBillNumberTooLong)
	c.optional(info.MobileNumber, maxMobileNumberLength, ErrMobileNumberTooLong)
	c.optional(info.Purpose, maxPurposeLength, ErrPurposeTooLong)
	c.check(func() error {
		return validateLanguageTemplate(info.AltLanguagePreference, info.AltMerchantName, info.AltMerchantCity)
	})
	c.check(func() error { return validateTimestamp(info.ExpirationTimestamp, info.Amount) })
	return c.err
}

func (info *MerchantInfo) validate() error {
	var c checker
	c.check(func() error { return validateUPIForGenerate(info.UPIAccountInfo, info.Currency) })
	c.check(func() error { return validateAccountID(info.BakongAccountID) })
	c.check(func() error { return validateCurrency(info.Currency) })
	c.check(func() error { return validateAmount(info.Amount, info.Currency) })
	c.check(func() error { return validateMerchantID(info.MerchantID) })
	c.check(func() error { return validateAcquiringBank(info.AcquiringBank) })
	c.check(func() error { return validateMerchantCategoryCode(info.MerchantCategoryCode) })
	c.check(func() error { return validateMerchantName(info.MerchantName) })
	c.check(func() error { return validateMerchantCity(info.MerchantCity) })
	c.optional(info.TerminalLabel, maxTerminalLabelLength, ErrTerminalLabelTooLong)
	c.optional(info.StoreLabel, maxStoreLabelLength, ErrStoreLabelTooLong)
	c.optional(info.BillNumber, maxBillNumberLength, ErrBillNumberTooLong)
	c.optional(info.MobileNumber, maxMobileNumberLength, ErrMobileNumberTooLong)
	c.optional(info.Purpose, maxPurposeLength, ErrPurposeTooLong)
	c.check(func() error {
		return validateLanguageTemplate(info.AltLanguagePreference, info.AltMerchantName, info.AltMerchantCity)
	})
	c.check(func() error { return validateTimestamp(info.ExpirationTimestamp, info.Amount) })
	return c.err
}

func (data *DecodedData) validate() error {
	var c checker
	var isDynamic bool
	c.check(func() error { return validateCRC(data.CRC) })
	c.check(func() error { return validatePayloadFormatIndicator(data.PayloadFormatIndicator) })
	c.check(func() error {
		var err error
		isDynamic, err = validatePointOfInitiationMethod(data.PointOfInitiationMethod)
		return err
	})
	c.check(func() error { return validateMerchantType(data.MerchantType) })
	c.check(func() error { return validateAccountID(data.BakongAccountID) })
	c.check(func() error {
		if data.MerchantType == Merchant {
			return validateMerchantID(data.MerchantID)
		}
		return nil
	})
	c.optional(data.AccountInfo, maxAccountIDLength, ErrAccountInfoTooLong)
	c.optional(data.AcquiringBank, maxAcquiringBankLength, ErrAcquiringBankTooLong)
	c.check(func() error { return validateDecodedMerchantCategoryCode(data.MerchantCategoryCode) })
	c.check(func() error { return validateTransactionCurrency(data.TransactionCurrency) })
	c.check(func() error { return data.validateTransactionAmount() })
	c.check(func() error { return validateCountryCode(data.CountryCode) })
	c.check(func() error { return validateMerchantName(data.MerchantName) })
	c.check(func() error { return validateMerchantCity(data.MerchantCity) })
	c.optional(data.BillNumber, maxBillNumberLength, ErrBillNumberTooLong)
	c.optional(data.MobileNumber, maxMobileNumberLength, ErrMobileNumberTooLong)
	c.optional(data.StoreLabel, maxStoreLabelLength, ErrStoreLabelTooLong)
	c.optional(data.TerminalLabel, maxTerminalLabelLength, ErrTerminalLabelTooLong)
	c.optional(data.Purpose, maxPurposeLength, ErrPurposeTooLong)
	c.check(func() error {
		return validateUPIForDecode(data.UPIAccountInfo, data.TransactionCurrency, data.CountryCode)
	})
	if c.err != nil {
		return c.err
	}
	if isDynamic {
		return data.validateDynamicFields()
	}
	return nil
}

func (data *DecodedData) validateTransactionAmount() error {
	if data.TransactionAmount == "" {
		return nil
	}
	amount, err := strconv.ParseFloat(data.TransactionAmount, 64)
	if err != nil || amount < 0 {
		return ErrInvalidAmount
	}
	if amount == 0 {
		return nil
	}
	formatted := strconv.FormatFloat(amount, 'f', -1, 64)
	if len(formatted) > maxAmountLength {
		return ErrInvalidAmount
	}
	switch data.TransactionCurrency {
	case khrCode:
		if math.Abs(amount-math.Floor(amount)) > 1e-9 { //nolint:mnd // float epsilon
			return ErrInvalidAmount
		}
	case usdCode:
		rounded := math.Round(amount*100) / 100 //nolint:mnd // cents multiplier
		if math.Abs(amount-rounded) > 1e-9 {    //nolint:mnd // float epsilon
			return ErrInvalidAmount
		}
	}
	return nil
}

func (data *DecodedData) validateDynamicFields() error {
	if strings.TrimSpace(data.TransactionAmount) == "" {
		return ErrInvalidDynamicKHQR
	}
	if data.ExpirationTimestamp == "" {
		return ErrExpirationRequired
	}
	if len(data.ExpirationTimestamp) != 13 { //nolint:mnd // timestamp string length
		return ErrInvalidTimestamp
	}
	ts, err := strconv.ParseInt(data.ExpirationTimestamp, 10, 64)
	if err != nil {
		return ErrInvalidTimestamp
	}
	if time.Now().UnixMilli() > ts {
		return ErrKHQRExpired
	}
	return nil
}
