package khqr

import "fmt"

// Error represents a KHQR validation or processing error.
type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("khqr: %s (code %d)", e.Message, e.Code)
}

// Is supports errors.Is by comparing error codes.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// Predefined errors for KHQR validation and processing.
var (
	ErrAccountIDRequired    = &Error{Code: 1, Message: "Bakong Account ID cannot be null or empty"}
	ErrMerchantNameRequired = &Error{Code: 2, Message: "Merchant name cannot be null or empty"}
	ErrAccountIDInvalid     = &Error{Code: 3, Message: "Bakong Account ID is invalid"}
	ErrInvalidAmount        = &Error{Code: 4, Message: "Amount is invalid"}
	ErrMerchantTypeRequired = &Error{Code: 5, Message: "Merchant type cannot be null or empty"}
	ErrAccountIDTooLong     = &Error{Code: 6, Message: "Bakong Account ID Length is invalid"}
	ErrMerchantNameTooLong  = &Error{Code: 7, Message: "Merchant Name Length is invalid"}
	ErrInvalidQR            = &Error{Code: 8, Message: "KHQR provided is invalid"}
	ErrCurrencyRequired     = &Error{Code: 9, Message: "Currency type cannot be null or empty"}
	ErrBillNumberTooLong    = &Error{Code: 10, Message: "Bill Number Length is invalid"}
	ErrStoreLabelTooLong    = &Error{Code: 11, Message: "Store Label Length is invalid"}
	ErrTerminalLabelTooLong = &Error{Code: 12, Message: "Terminal Label Length is invalid"}
	// ErrConnectionTimeout              = &Error{Code: 13, Message: "Cannot reach Bakong Open API service. Please check internet connection"}
	// ErrInvalidDeepLinkSourceInfo      = &Error{Code: 14, Message: "Source Info for Deep Link is invalid"}
	// ErrInternalServerError            = &Error{Code: 15, Message: "Internal server error"}
	ErrPayloadFormatIndicatorTooLong  = &Error{Code: 16, Message: "Payload Format Indicator Length is invalid"}
	ErrPointOfInitiationMethodTooLong = &Error{Code: 17, Message: "Point of Initiation Length is invalid"}
	ErrMerchantCategoryCodeTooLong    = &Error{Code: 18, Message: "Merchant Category Length is invalid"}
	ErrTransactionCurrencyTooLong     = &Error{Code: 19, Message: "Transaction Currency Length is invalid"}
	ErrCountryCodeTooLong             = &Error{Code: 20, Message: "Country Code Length is invalid"}
	ErrMerchantCityTooLong            = &Error{Code: 21, Message: "Merchant City Length is invalid"}
	ErrCRCInvalid                     = &Error{Code: 22, Message: "CRC Length is invalid"}
	ErrPayloadFormatIndicatorRequired = &Error{Code: 23, Message: "Payload Format Indicator cannot be null or empty"}
	ErrCRCRequired                    = &Error{Code: 24, Message: "CRC cannot be null or empty"}
	ErrMerchantCategoryCodeRequired   = &Error{Code: 25, Message: "Merchant Category cannot be null or empty"}
	ErrCountryCodeRequired            = &Error{Code: 26, Message: "Country Code cannot be null or empty"}
	ErrMerchantCityRequired           = &Error{Code: 27, Message: "Merchant City cannot be null or empty"}
	ErrInvalidCurrency                = &Error{Code: 28, Message: "Unsupported currency"}
	// ErrInvalidDeepLinkURL             = &Error{Code: 29, Message: "Deep Link URL is not valid"}
	ErrMerchantIDRequired    = &Error{Code: 30, Message: "Merchant ID cannot be null or empty"}
	ErrAcquiringBankRequired = &Error{Code: 31, Message: "Acquiring Bank cannot be null or empty"}
	ErrMerchantIDTooLong     = &Error{Code: 32, Message: "Merchant ID Length is invalid"}
	ErrAcquiringBankTooLong  = &Error{Code: 33, Message: "Acquiring Bank Length is invalid"}
	ErrMobileNumberTooLong   = &Error{Code: 34, Message: "Mobile Number Length is invalid"}
	// ErrTagNotInOrder                  = &Error{Code: 35, Message: "Tag is not in order"}
	ErrAccountInfoTooLong             = &Error{Code: 36, Message: "Account Information Length is invalid"}
	ErrLanguagePreferenceRequired     = &Error{Code: 37, Message: "Language Preference cannot be null or empty"}
	ErrLanguagePreferenceTooLong      = &Error{Code: 38, Message: "Language Preference Length is invalid"}
	ErrMerchantNameAltRequired        = &Error{Code: 39, Message: "Merchant Name Alternate Language cannot be null or empty"}
	ErrMerchantNameAltTooLong         = &Error{Code: 40, Message: "Merchant Name Alternate Language Length is invalid"}
	ErrMerchantCityAltTooLong         = &Error{Code: 41, Message: "Merchant City Alternate Language Length is invalid"}
	ErrPurposeTooLong                 = &Error{Code: 42, Message: "Purpose of Transaction Length is invalid"}
	ErrUPITooLong                     = &Error{Code: 43, Message: "Upi Account Information Length is invalid"}
	ErrUPINotSupportUSD               = &Error{Code: 44, Message: "KHQR does not support UPI Account Information with USD currency"}
	ErrExpirationRequired             = &Error{Code: 45, Message: "Expiration timestamp is required for dynamic KHQR"}
	ErrKHQRExpired                    = &Error{Code: 46, Message: "This dynamic KHQR has expired"}
	ErrInvalidDynamicKHQR             = &Error{Code: 47, Message: "This dynamic KHQR has invalid field transaction amount"}
	ErrPointOfInitiationMethodInvalid = &Error{Code: 48, Message: "Point of Initiation Method is invalid"}
	ErrInvalidTimestamp               = &Error{Code: 49, Message: "Expiration timestamp length is invalid"}
	ErrExpirationInPast               = &Error{Code: 50, Message: "Expiration timestamp is in the past"}
	ErrMerchantCategoryCodeInvalid    = &Error{Code: 51, Message: "Invalid Merchant Category Code"}
)
