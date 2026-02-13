package khqr

import (
	"crypto/md5" //nolint:gosec // MD5 used for KHQR SDK compatibility, not security
	"fmt"
)

// IndividualInfo contains information for generating an individual KHQR code.
type IndividualInfo struct {
	BakongAccountID       string
	Currency              Currency // default KHR
	Amount                float64  // 0 means no amount (static QR)
	MerchantName          string
	MerchantCity          string
	AcquiringBank         string // optional
	AccountInfo           string // optional
	UPIAccountInfo        string // optional
	BillNumber            string // optional
	StoreLabel            string // optional
	TerminalLabel         string // optional
	MobileNumber          string // optional
	Purpose               string // optional
	AltLanguagePreference string // optional
	AltMerchantName       string // optional
	AltMerchantCity       string // optional
	ExpirationTimestamp   int64  // unix ms, required if Amount > 0
	MerchantCategoryCode  string // optional, defaults to "5999"
}

// MerchantInfo contains information for generating a merchant KHQR code.
type MerchantInfo struct {
	BakongAccountID       string
	Currency              Currency // default KHR
	Amount                float64  // 0 means no amount (static QR)
	MerchantName          string
	MerchantCity          string
	MerchantID            string // required
	AcquiringBank         string // required
	UPIAccountInfo        string // optional
	BillNumber            string // optional
	StoreLabel            string // optional
	TerminalLabel         string // optional
	MobileNumber          string // optional
	Purpose               string // optional
	AltLanguagePreference string // optional
	AltMerchantName       string // optional
	AltMerchantCity       string // optional
	ExpirationTimestamp   int64  // unix ms, required if Amount > 0
	MerchantCategoryCode  string // optional, defaults to "5999"
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
