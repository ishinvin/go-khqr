package khqr

import (
	"crypto/md5" //nolint:gosec // MD5 used for KHQR SDK compatibility, not security
	"fmt"
)

// IndividualInfo contains information for generating an individual KHQR code.
type IndividualInfo struct {
	// Required fields.
	BakongAccountID string // must contain "@" (e.g. "user@bank")
	MerchantName    string // max 25 characters

	// Optional fields with defaults.
	Currency             Currency // defaults to KHR
	MerchantCity         string   // defaults to "Phnom Penh", max 15 characters
	MerchantCategoryCode string   // defaults to "5999"

	// Optional fields (dynamic QR).
	Amount              float64 // 0 means static QR; KHR must be whole, USD up to 2 decimals
	ExpirationTimestamp int64   // unix ms, required when Amount > 0

	// Optional fields.
	AcquiringBank         string // max 32 characters
	AccountInfo           string // max 32 characters
	UPIAccountInfo        string // not supported with USD, max 99 characters
	BillNumber            string // max 25 characters
	StoreLabel            string // max 25 characters
	TerminalLabel         string // max 25 characters
	MobileNumber          string // max 25 characters
	Purpose               string // max 25 characters
	AltLanguagePreference string // ISO 639-1 (2 chars); requires AltMerchantName
	AltMerchantName       string // required when AltLanguagePreference is set, max 25 characters
	AltMerchantCity       string // max 15 characters
}

// MerchantInfo contains information for generating a merchant KHQR code.
type MerchantInfo struct {
	// Required fields.
	BakongAccountID string // must contain "@" (e.g. "merchant@bank")
	MerchantName    string // max 25 characters
	MerchantCity    string // max 15 characters
	MerchantID      string // max 32 characters
	AcquiringBank   string // max 32 characters

	// Optional fields with defaults.
	Currency             Currency // defaults to KHR
	MerchantCategoryCode string   // defaults to "5999"

	// Optional fields (dynamic QR).
	Amount              float64 // 0 means static QR; KHR must be whole, USD up to 2 decimals
	ExpirationTimestamp int64   // unix ms, required when Amount > 0

	// Optional fields.
	UPIAccountInfo        string // not supported with USD, max 99 characters
	BillNumber            string // max 25 characters
	StoreLabel            string // max 25 characters
	TerminalLabel         string // max 25 characters
	MobileNumber          string // max 25 characters
	Purpose               string // max 25 characters
	AltLanguagePreference string // ISO 639-1 (2 chars); requires AltMerchantName
	AltMerchantName       string // required when AltLanguagePreference is set, max 25 characters
	AltMerchantCity       string // max 15 characters
}

// Data contains the generated QR string.
type Data struct {
	QR string
}

// String returns the QR payload string.
func (d *Data) String() string {
	return d.QR
}

// MD5 returns the hex-encoded MD5 hash of the QR string.
func (d *Data) MD5() string {
	hash := md5.Sum([]byte(d.QR)) //nolint:gosec // MD5 used for KHQR SDK compatibility, not security
	return fmt.Sprintf("%x", hash)
}

// DecodedData contains all decoded fields from a KHQR string.
type DecodedData struct {
	PayloadFormatIndicator  string
	PointOfInitiationMethod string
	BakongAccountID         string
	MerchantID              string
	AccountInfo             string
	AcquiringBank           string
	MerchantType            MerchantType
	TransactionCurrency     string
	MerchantName            string
	TransactionAmount       string
	MerchantCategoryCode    string
	CountryCode             string
	MerchantCity            string
	BillNumber              string
	StoreLabel              string
	TerminalLabel           string
	MobileNumber            string
	CreationTimestamp       string
	ExpirationTimestamp     string
	CRC                     string
	UPIAccountInfo          string
	Purpose                 string
	AltLanguagePreference   string
	AltMerchantName         string
	AltMerchantCity         string
}
