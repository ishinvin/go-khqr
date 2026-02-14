package khqr

import (
	"strings"
	"time"
)

// qrParams holds the type-agnostic inputs for QR generation.
// Callers are responsible for validation before constructing this.
type qrParams struct {
	accountTag   string // "29" or "30"
	accountValue string // pre-built subtag content

	merchantName   string
	merchantCity   string
	categoryCode   string   // already defaulted
	currency       Currency // already defaulted
	amount         float64
	expiration     int64
	upiAccountInfo string

	billNumber    string
	mobileNumber  string
	storeLabel    string
	terminalLabel string
	purpose       string

	altLanguagePreference string
	altMerchantName       string
	altMerchantCity       string
}

// generate builds a KHQR payload from type-agnostic parameters.
// Tags are written in ascending order per the EMV QR Code specification.
// Tags (00, 01, 15, 29/30, 52, 53, 54, 58, 59, 60, 62, 64, 99).
func generate(p *qrParams) *Data {
	isDynamic := p.amount > 0
	var b strings.Builder

	// Payload Format Indicator (tag 00)
	b.WriteString(encodeTLV(tagPayloadFormatIndicator, defaultPayloadFormatIndicator))

	// Point of Initiation (tag 01)
	poi := staticQR
	if isDynamic {
		poi = dynamicQR
	}
	b.WriteString(encodeTLV(tagPointOfInitiation, poi))

	// UnionPay (tag 15)
	if p.upiAccountInfo != "" {
		b.WriteString(encodeTLV(tagUnionPay, p.upiAccountInfo))
	}

	// Individual/Merchant Account (tag 29/30)
	b.WriteString(encodeTLV(p.accountTag, p.accountValue))

	// Merchant Category Code (tag 52)
	b.WriteString(encodeTLV(tagMerchantCategoryCode, p.categoryCode))

	// Transaction Currency (tag 53)
	b.WriteString(encodeTLV(tagCurrency, formatCurrency(p.currency)))

	// Transaction Amount (tag 54)
	if isDynamic {
		b.WriteString(encodeTLV(tagAmount, formatAmount(p.amount, p.currency)))
	}

	// Country Code (tag 58)
	b.WriteString(encodeTLV(tagCountryCode, defaultCountryCode))

	// Merchant Name (tag 59)
	b.WriteString(encodeTLV(tagMerchantName, p.merchantName))

	// Merchant City (tag 60)
	b.WriteString(encodeTLV(tagMerchantCity, p.merchantCity))

	// Additional Data (tag 62)
	if ad := buildAdditionalData(p.billNumber, p.mobileNumber, p.storeLabel, p.terminalLabel, p.purpose); ad != "" {
		b.WriteString(encodeTLV(tagAdditionalData, ad))
	}

	// Language Template (tag 64)
	if lt := buildLanguageTemplate(p.altLanguagePreference, p.altMerchantName, p.altMerchantCity); lt != "" {
		b.WriteString(encodeTLV(tagLanguageTemplate, lt))
	}

	// Timestamp (tag 99) — dynamic only
	if isDynamic {
		var ts strings.Builder
		ts.WriteString(encodeTLV(subtagCreationTimestamp, formatTimestamp(time.Now().UnixMilli())))
		ts.WriteString(encodeTLV(subtagExpirationTimestamp, formatTimestamp(p.expiration)))
		b.WriteString(encodeTLV(tagTimestamp, ts.String()))
	}

	// CRC (tag 63)
	b.WriteString(tagCRC + "04")
	payload := b.String()
	payload += crc16Hex(payload)

	return &Data{QR: payload}
}

// generateIndividual builds a KHQR payload for an individual payment.
func generateIndividual(info *IndividualInfo) (*Data, error) {
	if info.Currency == 0 {
		info.Currency = KHR
	}
	if info.MerchantCategoryCode == "" {
		info.MerchantCategoryCode = defaultMerchantCategoryCode
	}
	if info.MerchantCity == "" {
		info.MerchantCity = defaultMerchantCity
	}

	if err := info.validate(); err != nil {
		return nil, err
	}

	// Build individual account subtags (tag 29)
	var acc tlvWriter
	acc.WriteString(encodeTLV(subtagGlobalID, info.BakongAccountID))
	acc.writeTLV(subtagAccountInfo, info.AccountInfo)
	acc.writeTLV(subtagAcquiringBank, info.AcquiringBank)

	return generate(&qrParams{
		accountTag:            tagIndividualAccount,
		accountValue:          acc.String(),
		merchantName:          info.MerchantName,
		merchantCity:          info.MerchantCity,
		categoryCode:          info.MerchantCategoryCode,
		currency:              info.Currency,
		amount:                info.Amount,
		expiration:            info.ExpirationTimestamp,
		upiAccountInfo:        info.UPIAccountInfo,
		billNumber:            info.BillNumber,
		mobileNumber:          info.MobileNumber,
		storeLabel:            info.StoreLabel,
		terminalLabel:         info.TerminalLabel,
		purpose:               info.Purpose,
		altLanguagePreference: info.AltLanguagePreference,
		altMerchantName:       info.AltMerchantName,
		altMerchantCity:       info.AltMerchantCity,
	}), nil
}

// generateMerchant builds a KHQR payload for a merchant payment.
func generateMerchant(info *MerchantInfo) (*Data, error) {
	if info.Currency == 0 {
		info.Currency = KHR
	}
	if info.MerchantCategoryCode == "" {
		info.MerchantCategoryCode = defaultMerchantCategoryCode
	}

	if err := info.validate(); err != nil {
		return nil, err
	}

	// Build merchant account subtags (tag 30)
	var acc strings.Builder
	acc.WriteString(encodeTLV(subtagGlobalID, info.BakongAccountID))
	acc.WriteString(encodeTLV(subtagMerchantID, info.MerchantID))
	acc.WriteString(encodeTLV(subtagAcquiringBank, info.AcquiringBank))

	return generate(&qrParams{
		accountTag:            tagMerchantAccount,
		accountValue:          acc.String(),
		merchantName:          info.MerchantName,
		merchantCity:          info.MerchantCity,
		categoryCode:          info.MerchantCategoryCode,
		currency:              info.Currency,
		amount:                info.Amount,
		expiration:            info.ExpirationTimestamp,
		upiAccountInfo:        info.UPIAccountInfo,
		billNumber:            info.BillNumber,
		mobileNumber:          info.MobileNumber,
		storeLabel:            info.StoreLabel,
		terminalLabel:         info.TerminalLabel,
		purpose:               info.Purpose,
		altLanguagePreference: info.AltLanguagePreference,
		altMerchantName:       info.AltMerchantName,
		altMerchantCity:       info.AltMerchantCity,
	}), nil
}

// buildAdditionalData constructs the additional data field (tag 62) content.
func buildAdditionalData(billNumber, mobileNumber, storeLabel, terminalLabel, purpose string) string {
	var b tlvWriter
	b.writeTLV(subtagBillNumber, billNumber)
	b.writeTLV(subtagMobileNumber, mobileNumber)
	b.writeTLV(subtagStoreLabel, storeLabel)
	b.writeTLV(subtagTerminalLabel, terminalLabel)
	b.writeTLV(subtagPurpose, purpose)
	return b.String()
}

// buildLanguageTemplate constructs the language template field (tag 64) content.
func buildLanguageTemplate(preference, nameAlt, cityAlt string) string {
	if preference == "" {
		return ""
	}
	var b tlvWriter
	b.WriteString(encodeTLV(subtagLanguagePreference, preference))
	b.writeTLV(subtagMerchantNameAlt, nameAlt)
	b.writeTLV(subtagMerchantCityAlt, cityAlt)
	return b.String()
}
