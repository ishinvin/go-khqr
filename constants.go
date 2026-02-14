package khqr

import "fmt"

// Currency represents ISO 4217 numeric currency codes supported by KHQR.
type Currency int

const (
	KHR Currency = 116
	USD Currency = 840
)

// String returns the ISO 4217 alphabetic code for the currency.
func (c Currency) String() string {
	switch c {
	case KHR:
		return "KHR"
	case USD:
		return "USD"
	default:
		return fmt.Sprintf("Currency(%d)", int(c))
	}
}

// MerchantType indicates whether a KHQR code is for an individual or merchant.
type MerchantType string

const (
	Individual MerchantType = "individual"
	Merchant   MerchantType = "merchant"
)

// EMV tag codes (internal — users don't construct TLV manually)
const (
	tagPayloadFormatIndicator = "00"
	tagPointOfInitiation      = "01"
	tagUnionPay               = "15"
	tagIndividualAccount      = "29"
	tagMerchantAccount        = "30"
	tagMerchantCategoryCode   = "52"
	tagCurrency               = "53"
	tagAmount                 = "54"
	tagCountryCode            = "58"
	tagMerchantName           = "59"
	tagMerchantCity           = "60"
	tagAdditionalData         = "62"
	tagCRC                    = "63"
	tagLanguageTemplate       = "64"
	tagTimestamp              = "99"
)

// Subtag codes for merchant account tags (29/30)
const (
	subtagGlobalID      = "00"
	subtagMerchantID    = "01"
	subtagAccountInfo   = "01"
	subtagAcquiringBank = "02"
)

// Subtag codes for additional data (tag 62)
const (
	subtagBillNumber    = "01"
	subtagMobileNumber  = "02"
	subtagStoreLabel    = "03"
	subtagTerminalLabel = "07"
	subtagPurpose       = "08"
)

// Subtag codes for language template (tag 64)
const (
	subtagLanguagePreference = "00"
	subtagMerchantNameAlt    = "01"
	subtagMerchantCityAlt    = "02"
)

// Subtag codes for timestamp (tag 99)
const (
	subtagCreationTimestamp   = "00"
	subtagExpirationTimestamp = "01"
)

// Internal default values
const (
	defaultPayloadFormatIndicator = "01"
	defaultMerchantCategoryCode   = "5999"
	defaultMerchantCity           = "Phnom Penh"
	defaultCountryCode            = "KH"
	staticQR                      = "11"
	dynamicQR                     = "12"
)

// Maximum field lengths
const (
	maxAccountIDLength       = 32
	maxMerchantNameLength    = 25
	maxMerchantCityLength    = 15
	maxAmountLength          = 13
	maxBillNumberLength      = 25
	maxStoreLabelLength      = 25
	maxTerminalLabelLength   = 25
	maxMobileNumberLength    = 25
	maxPurposeLength         = 25
	maxMerchantIDLength      = 32
	maxAcquiringBankLength   = 32
	maxUPILength             = 99
	maxMerchantNameAltLength = 25
	maxMerchantCityAltLength = 15
	languagePreferenceLength = 2
)
