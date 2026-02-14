package khqr

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// mockTimestampAndCRC appends tag 99 (creation + expiration timestamps)
// and tag 63 (CRC) to a partial KHQR string for testing dynamic QR codes.
func mockTimestampAndCRC(base string) string {
	var b strings.Builder
	b.WriteString(base)

	var ts strings.Builder
	ts.WriteString(encodeTLV(subtagCreationTimestamp, formatTimestamp(time.Now().UnixMilli())))
	ts.WriteString(encodeTLV(subtagExpirationTimestamp, formatTimestamp(time.Now().Add(5*time.Minute).UnixMilli())))
	b.WriteString(encodeTLV(tagTimestamp, ts.String()))

	b.WriteString(tagCRC + "04")
	payload := b.String()
	payload += crc16Hex(payload)
	return payload
}

func TestVerifyValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		qr   string
	}{
		{
			"individual_usd_khmer_name",
			mockTimestampAndCRC("00020101021229180014jonhsmith@nbcq52045999530384054031.05802KH5906កូរូណា6010Phnom Penh"),
		},
		{
			"individual_khr_with_mobile",
			mockTimestampAndCRC("00020101021229180014jonhsmith@nbcq520459995303116540750000.05802KH5910Jonh Smith6010Phnom Penh6215021185512345678"),
		},
		{
			"individual_usd_with_bill",
			mockTimestampAndCRC("00020101021229180014jonhsmith@nbcq52045999530384054031.05802KH5910Jonh Smith6010PHNOM PENH62210117INV-2021-07-65822"),
		},
		{
			"merchant_usd_with_additional",
			mockTimestampAndCRC("00020101021230400014jonhsmith@devb01061234560208Dev Bank520459995303840540410.05802KH5910Jonh Smith6010Phnom Penh62360117INV-2021-07-65822021185512345678"),
		},
		{
			"individual_tags_reversed",
			mockTimestampAndCRC("62210117INV-2021-07-658226010PHNOM PENH5910Jonh Smith5802KH5401153038405204599929180014jonhsmith@nbcq010212000201"),
		},
		{
			"individual_tags_reversed_with_additional",
			mockTimestampAndCRC("6250070201030412340211855854989940117INV-2021-07-658226010PHNOM PENH5910Jonh Smith5802KH540115303840520459992926070412340014jonhsmith@nbcq010212000201"),
		},
		// Static QR codes with pre-computed CRC
		{
			"static_with_timestamps",
			"00020101021129180014jonhsmith@nbcq5204599953031165802KH5910Jonh Smith6010Phnom Penh62290225012345678901234567890123499170013167766094781963040CC8",
		},
		{
			"static_with_language_template",
			"00020101021129180014jonhsmith@nbcq5204599953031165802KH5910Jonh Smith6009Siam Reap64280002km0108ចន ស្មីន0206សៀមរាប99170013168733740175863047488",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := Verify(tt.qr)
			if err != nil {
				t.Errorf("Verify() = %v, want nil", err)
			}
		})
	}
}

func TestVerifyError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		qr      string
		wantErr error
	}{
		// Invalid QR (too short or unparseable)
		{"empty_string", "", ErrInvalidQR},
		{"too_short", "ABC", ErrInvalidQR},
		{"too_short_with_spaces", "ABC         ", ErrInvalidQR},
		{"too_short_hex", "EAB3", ErrInvalidQR},

		// CRC errors (format or checksum mismatch)
		{"khmer_chars_no_crc", "កខគឃង", ErrCRCInvalid},
		{"crc_checksum_mismatch", "00020101021126150011jonhsmith@nbcq5204599953031165802KH5912Jonh Smith6010Phnom Penh63046D28", ErrCRCInvalid},
		{"crc_too_short", "00020101021126150011jonhsmith@nbcq5204599953031165802KH5912Jonh Smith6010Phnom Penh63046D", ErrCRCInvalid},
		{"missing_crc_tag", "00020101021230190015", ErrCRCInvalid},
		{"crc_mismatch_with_prefix", "10263041234", ErrCRCInvalid},
		{"crc_mismatch_minimal", "63041234", ErrCRCInvalid},

		// Merchant type required
		{"merchant_type_required_tag26", "00020101021126200016coffeeklang@pras63045927", ErrMerchantTypeRequired},
		{"merchant_type_required_no_account_tag", "00020101021252045999530384054035.05802KH5916john smith actor6010Phnom Penh62150111Invoice#06999170013161343857589263042494", ErrMerchantTypeRequired},

		// Payload format indicator
		{"pfi_required_only_crc", "63046007", ErrPayloadFormatIndicatorRequired},
		{"pfi_required_no_tag00", "01021230190015john_smith@devb52045999530384054035.05802KH5916john smith actor6010Phnom Penh62150111Invoice#06999170013161343857589263044598", ErrPayloadFormatIndicatorRequired},
		{"pfi_too_long", "000301001021230190015john_smith@devb630493FD", ErrPayloadFormatIndicatorTooLong},

		// Point of initiation method
		{"poi_too_long", "000201010312330190015john_smith@devb63040FC8", ErrPointOfInitiationMethodTooLong},
		{"poi_invalid", "00020101021329180014vandoeurn@devb52045999530384054031015802KH5909Vandoeurn6010Phnom Penh993400131727249364974011317272494849746304DCEF", ErrPointOfInitiationMethodInvalid},

		// Bakong account ID
		{"account_id_invalid", "000201010212301900151234567890123456304D807", ErrAccountIDInvalid},
		{"account_id_too_long", "00020101021230440040123456789012345678901234567890123456789063041747", ErrAccountIDTooLong},

		// Merchant category code
		{"mcc_too_long", "00020101021230410015john_smith@devb01061234560208Dev Bank52055999153038405405100.05802KH5910John Smith6010PHNOM PENH62530106#123450211855122334550311Coffee Shop0709Cashier_19917001316262478678616304954C", ErrMerchantCategoryCodeTooLong},
		{"mcc_required", "00020101021229190015john_smith@devb530384054035.05802KH5916john smith actor6010Phnom Penh62150111Invoice#0699917001316134385758926304F926", ErrMerchantCategoryCodeRequired},
		{"mcc_invalid_hex", "00020101021129180014jonhsmith@nbcq52041A2B53031165802KH5910Jonh Smith6010Phnom Penh6304F7FD", ErrMerchantCategoryCodeInvalid},
		{"mcc_invalid_negative", "00020101021129180014jonhsmith@nbcq5204-10053031165802KH5910Jonh Smith6010Phnom Penh6304038A", ErrMerchantCategoryCodeInvalid},

		// Transaction currency
		{"currency_too_long", "00020101021229190015john_smith@devb52045999530488406304FEA6", ErrTransactionCurrencyTooLong},
		{"currency_required", "00020101021229190015john_smith@devb5204599954035.05802KH5916john smith actor6010Phnom Penh62150111Invoice#0699917001316134385758926304E7DF", ErrCurrencyRequired},
		{"currency_unsupported", "00020101021229190015john_smith@devb52045999530384954035.05802KH5916john smith actor6010Phnom Penh9917001316134385758926304951E", ErrInvalidCurrency},

		// Country code
		{"country_code_too_long", "00020101021229190015john_smith@devb52045999530384054035.05803KKH63042FD2", ErrCountryCodeTooLong},
		{"country_code_required", "00020101021229190015john_smith@devb52045999530384054035.05916john smith actor6010Phnom Penh62150111Invoice#06999170013161343857589263040023", ErrCountryCodeRequired},

		// Merchant name
		{"merchant_name_too_long", "00020101021229190015john_smith@devb52045999530384054035.05802KH59307PL7EvxHpgpP4jT4uMgegaYqgv3Ehb6010Phnom Penh6304399D", ErrMerchantNameTooLong},
		{"merchant_name_required_empty", "00020101021229190015john_smith@devb52045999530384054035.05802KH59006010Phnom Penh630414A7", ErrMerchantNameRequired},
		{"merchant_name_required_missing", "00020101021229190015john_smith@devb52045999530384054035.05802KH6010Phnom Penh62150111Invoice#06999170013161343857589263043518", ErrMerchantNameRequired},

		// Merchant city
		{"merchant_city_too_long", "00020101021229190015john_smith@devb52045999530384054035.05802KH5916john smith actor60170123456789012345663040EF3", ErrMerchantCityTooLong},
		{"merchant_city_required_empty", "00020101021229190015john_smith@devb52045999530384054035.05802KH5916john smith actor600062150111Invoice#0699917001316134385758926304389D", ErrMerchantCityRequired},
		{"merchant_city_required_missing", "00020101021229190015john_smith@devb52045999530384054035.05802KH5916john smith actor62150111Invoice#06999170013161343857589263043B26", ErrMerchantCityRequired},

		// Additional data field lengths (merchant QR)
		{"bill_number_too_long", "00020101021230410015john_smith@devb01061234560208Dev Bank5204599953038405405100.05802KH5910John Smith6010PHNOM PENH62740127#123456789012345678901234560211855122334550311Coffee Shop0709Cashier_199170013162624810862363040A18", ErrBillNumberTooLong},
		{"store_label_too_long", "00020101021230410015john_smith@devb01061234560208Dev Bank5204599953038405405100.05802KH5910John Smith6010PHNOM PENH62730106#123450211855122334550331Coffee Shop123456789012345678900709Cashier_199170013162624818838863046A71", ErrStoreLabelTooLong},
		{"terminal_label_too_long", "00020101021230410015john_smith@devb01061234560208Dev Bank5204599953038405405100.05802KH5910John Smith6010PHNOM PENH62760106#123450211855122334550311Coffee Shop0732Cashier_1234567890123456789012349917001316262482295746304F5E3", ErrTerminalLabelTooLong},

		// Account info / acquiring bank length
		{"account_info_too_long", "00020101021229670014jonhsmith@nbcq01331234567890123456789012345678901230208Dev Bank520459995303116540750000.05802KH5910Jonh Smith6009Siam Reap62550117INV-2021-07-658220211855123456780305BKK-107060123459917001316304651886946304E82B", ErrAccountInfoTooLong},
		{"acquiring_bank_too_long", "00020101021229550014jonhsmith@nbcq0233123456789012345678901234567890123520459995303116540750000.05802KH5910Jonh Smith6009Siam Reap62550117INV-2021-07-658220211855123456780305BKK-107060123459917001316304670122756304AE3F", ErrAcquiringBankTooLong},

		// UPI with USD
		{"upi_not_support_usd", "0002010102111531197401160052044645410994011112530390010ii_rr@devb01090000123460208Dev Bank520459995303840540105802TH5906សួស្ដី6002PP621502118551234567899170013171221890231163042A5A", ErrUPINotSupportUSD},

		// Dynamic QR errors
		{"invalid_dynamic_khqr_no_amount", "00020101021229180014vandoeurn@devb5204599953038405802KH5909Vandoeurn6010Phnom Penh993400131727249364974011317272494849746304B6F4", ErrInvalidDynamicKHQR},
		{"khqr_expired", "00020101021229180014vandoeurn@devb52045999530384054031015802KH5909Vandoeurn6010Phnom Penh9934001317272493649740113172724948497463047BD0", ErrKHQRExpired},
		{"expiration_required_no_tag99", "00020101021215311974011600520446ACLB1000231208129200016chantha_dev@ftcc520459995303116540410005802KH5911Chantha Dev6010Phnom Penh63046011", ErrExpirationRequired},
		{"expiration_required_creation_only", "00020101021215311974011600520446ACLB1000231208129200016chantha_dev@ftcc52045999530311654034005802KH5911Chantha Dev6010Phnom Penh99170013172974008531463046EEC", ErrExpirationRequired},
		{"expiration_timestamp_length_invalid", "00020101021229220018sopheak_leng2@ftcc520459995303116540410005802KH5912sopheak leng6010Phnom Penh9936001411732270085311011417323270145310630445D7", ErrInvalidTimestamp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := Verify(tt.qr)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Verify(%q) = %v, want %v", tt.qr, err, tt.wantErr)
			}
		})
	}
}
