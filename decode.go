package khqr

import "strings"

// fieldMaps holds pre-built tag -> field pointer maps used during decoding.
type fieldMaps struct {
	topLevel   map[string]*string
	individual map[string]*string
	merchant   map[string]*string
	additional map[string]*string
	language   map[string]*string
	timestamp  map[string]*string
}

// decodeMaps builds all field maps once for the decode pass.
func (data *DecodedData) decodeMaps() fieldMaps {
	return fieldMaps{
		topLevel: map[string]*string{
			tagPayloadFormatIndicator: &data.PayloadFormatIndicator,
			tagPointOfInitiation:      &data.PointOfInitiationMethod,
			tagUnionPay:               &data.UPIAccountInfo,
			tagMerchantCategoryCode:   &data.MerchantCategoryCode,
			tagCurrency:               &data.TransactionCurrency,
			tagAmount:                 &data.TransactionAmount,
			tagCountryCode:            &data.CountryCode,
			tagMerchantName:           &data.MerchantName,
			tagMerchantCity:           &data.MerchantCity,
			tagCRC:                    &data.CRC,
		},
		individual: map[string]*string{
			subtagGlobalID:      &data.BakongAccountID,
			subtagAccountInfo:   &data.AccountInfo,
			subtagAcquiringBank: &data.AcquiringBank,
		},
		merchant: map[string]*string{
			subtagGlobalID:      &data.BakongAccountID,
			subtagMerchantID:    &data.MerchantID,
			subtagAcquiringBank: &data.AcquiringBank,
		},
		additional: map[string]*string{
			subtagBillNumber:    &data.BillNumber,
			subtagMobileNumber:  &data.MobileNumber,
			subtagStoreLabel:    &data.StoreLabel,
			subtagTerminalLabel: &data.TerminalLabel,
			subtagPurpose:       &data.Purpose,
		},
		language: map[string]*string{
			subtagLanguagePreference: &data.AltLanguagePreference,
			subtagMerchantNameAlt:    &data.AltMerchantName,
			subtagMerchantCityAlt:    &data.AltMerchantCity,
		},
		timestamp: map[string]*string{
			subtagCreationTimestamp:   &data.CreationTimestamp,
			subtagExpirationTimestamp: &data.ExpirationTimestamp,
		},
	}
}

// decode parses a KHQR string without validating CRC, returning the decoded data.
func decode(qr string) (*DecodedData, error) {
	qr = strings.TrimSpace(qr)

	entries, err := parseTLV(qr)
	if err != nil {
		return nil, ErrInvalidQR
	}

	data := &DecodedData{}
	dm := data.decodeMaps()
	for _, entry := range entries {
		if err := decodeEntry(entry, data, dm); err != nil {
			return nil, err
		}
	}

	return data, nil
}

// decodeEntry dispatches a single top-level TLV entry to the appropriate handler.
func decodeEntry(entry tlv, data *DecodedData, dm fieldMaps) error {
	if ptr := dm.topLevel[entry.Tag]; ptr != nil {
		*ptr = entry.Value
		return nil
	}

	switch entry.Tag {
	case tagIndividualAccount:
		data.MerchantType = Individual
		return decodeSubtags(entry.Value, dm.individual)
	case tagMerchantAccount:
		data.MerchantType = Merchant
		return decodeSubtags(entry.Value, dm.merchant)
	case tagAdditionalData:
		return decodeSubtags(entry.Value, dm.additional)
	case tagLanguageTemplate:
		return decodeSubtags(entry.Value, dm.language)
	case tagTimestamp:
		return decodeSubtags(entry.Value, dm.timestamp)
	}
	return nil
}

// decodeSubtags parses nested TLV data and assigns values to the mapped fields.
func decodeSubtags(value string, fields map[string]*string) error {
	entries, err := parseTLV(value)
	if err != nil {
		return ErrInvalidQR
	}
	for _, e := range entries {
		if ptr := fields[e.Tag]; ptr != nil {
			*ptr = e.Value
		}
	}
	return nil
}
