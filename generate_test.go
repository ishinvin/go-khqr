package khqr

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// --- Generate Individual Success ---

func TestGenerateIndividualSuccess(t *testing.T) {
	t.Parallel()

	futureTS := time.Now().Add(2 * time.Minute).UnixMilli()

	tests := []struct {
		name string
		info *IndividualInfo
		want string // full expected QR from Java SDK; dynamic suffix will be stripped
	}{
		{
			"USD_amount_1",
			&IndividualInfo{
				BakongAccountID:     "jonhsmith@nbcq",
				MerchantName:        "Jonh Smith",
				MerchantCity:        "PHNOM PENH",
				Currency:            USD,
				Amount:              1.0,
				BillNumber:          "INV-2021-07-65822",
				ExpirationTimestamp: futureTS,
			},
			"00020101021229180014jonhsmith@nbcq520459995303840540115802KH5910Jonh Smith6010PHNOM PENH62210117INV-2021-07-658229934001317267489142870113172674951428763045F5D",
		},
		{
			"KHR_amount_50000_mobile",
			&IndividualInfo{
				BakongAccountID:     "jonhsmith@nbcq",
				MerchantName:        "Jonh Smith",
				Currency:            KHR,
				Amount:              50000,
				MobileNumber:        "85512345678",
				ExpirationTimestamp: futureTS,
			},
			"00020101021229180014jonhsmith@nbcq5204599953031165405500005802KH5910Jonh Smith6010Phnom Penh6215021185512345678993400131726749117626011317267497176266304CA4B",
		},
		{
			"KHR_default_currency_store_label",
			&IndividualInfo{
				BakongAccountID:     "jonhsmith@nbcq",
				MerchantName:        "Jonh Smith",
				Amount:              50000,
				StoreLabel:          "BKK-1",
				ExpirationTimestamp: futureTS,
			},
			"00020101021229180014jonhsmith@nbcq5204599953031165405500005802KH5910Jonh Smith6010Phnom Penh62090305BKK-19934001317267491176330113172674971763363040E36",
		},
		{
			"KHR_all_additional_data",
			&IndividualInfo{
				BakongAccountID:     "jonhsmith@nbcq",
				MerchantName:        "Jonh Smith",
				MerchantCity:        "Siam Reap",
				Currency:            KHR,
				Amount:              50000,
				BillNumber:          "INV-2021-07-65822",
				MobileNumber:        "85512345678",
				StoreLabel:          "BKK-1",
				TerminalLabel:       "012345",
				ExpirationTimestamp: futureTS,
			},
			"00020101021229180014jonhsmith@nbcq5204599953031165405500005802KH5910Jonh Smith6009Siam Reap62550117INV-2021-07-658220211855123456780305BKK-107060123459934001317267491176360113172674971763663049D4A",
		},
		{
			"KHR_with_account_info_and_bank",
			&IndividualInfo{
				BakongAccountID:     "jonhsmith@nbcq",
				AccountInfo:         "012345678",
				AcquiringBank:       "Dev Bank",
				MerchantName:        "Jonh Smith",
				MerchantCity:        "Siam Reap",
				Currency:            KHR,
				Amount:              50000,
				BillNumber:          "INV-2021-07-65822",
				MobileNumber:        "85512345678",
				StoreLabel:          "BKK-1",
				TerminalLabel:       "012345",
				ExpirationTimestamp: futureTS,
			},
			"00020101021229430014jonhsmith@nbcq01090123456780208Dev Bank5204599953031165405500005802KH5910Jonh Smith6009Siam Reap62550117INV-2021-07-658220211855123456780305BKK-107060123459934001317267491176380113172674971763863046CD1",
		},
		{
			"KHR_with_bank_only",
			&IndividualInfo{
				BakongAccountID:     "jonhsmith@nbcq",
				AcquiringBank:       "Dev Bank",
				MerchantName:        "Jonh Smith",
				MerchantCity:        "Siam Reap",
				Currency:            KHR,
				Amount:              50000,
				BillNumber:          "INV-2021-07-65822",
				MobileNumber:        "85512345678",
				StoreLabel:          "BKK-1",
				TerminalLabel:       "012345",
				ExpirationTimestamp: futureTS,
			},
			"00020101021229300014jonhsmith@nbcq0208Dev Bank5204599953031165405500005802KH5910Jonh Smith6009Siam Reap62550117INV-2021-07-658220211855123456780305BKK-107060123459934001317267491176400113172674971764063049E56",
		},
		{
			"KHR_with_account_info_only",
			&IndividualInfo{
				BakongAccountID:     "jonhsmith@nbcq",
				AccountInfo:         "012345678",
				MerchantName:        "Jonh Smith",
				MerchantCity:        "Siam Reap",
				Currency:            KHR,
				Amount:              50000,
				BillNumber:          "INV-2021-07-65822",
				MobileNumber:        "85512345678",
				StoreLabel:          "BKK-1",
				TerminalLabel:       "012345",
				ExpirationTimestamp: futureTS,
			},
			"00020101021229310014jonhsmith@nbcq01090123456785204599953031165405500005802KH5910Jonh Smith6009Siam Reap62550117INV-2021-07-658220211855123456780305BKK-10706012345993400131726749117643011317267497176436304220E",
		},
		{
			"KHR_amount_100",
			&IndividualInfo{
				BakongAccountID:     "jonhsmith@nbcq",
				AccountInfo:         "012345678",
				AcquiringBank:       "Dev Bank",
				MerchantName:        "Jonh Smith",
				MerchantCity:        "Siam Reap",
				Currency:            KHR,
				Amount:              100,
				BillNumber:          "INV-2021-07-65822",
				MobileNumber:        "85512345678",
				StoreLabel:          "BKK-1",
				TerminalLabel:       "012345",
				ExpirationTimestamp: futureTS,
			},
			"00020101021229430014jonhsmith@nbcq01090123456780208Dev Bank52045999530311654031005802KH5910Jonh Smith6009Siam Reap62550117INV-2021-07-658220211855123456780305BKK-10706012345993400131726749117644011317267497176446304F781",
		},
		{
			"KHR_khmer_merchant_name",
			&IndividualInfo{
				BakongAccountID:     "jonhsmith@nbcq",
				AccountInfo:         "012345678",
				AcquiringBank:       "Dev Bank",
				MerchantName:        "ចន ស្មីន",
				MerchantCity:        "Siam Reap",
				Currency:            KHR,
				Amount:              100,
				BillNumber:          "INV-2021-07-65822",
				MobileNumber:        "85512345678",
				StoreLabel:          "BKK-1",
				TerminalLabel:       "012345",
				ExpirationTimestamp: futureTS,
			},
			"00020101021229430014jonhsmith@nbcq01090123456780208Dev Bank52045999530311654031005802KH5908ចន ស្មីន6009Siam Reap62550117INV-2021-07-658220211855123456780305BKK-107060123459934001317267491176470113172674971764763044F9A",
		},
		{
			"static_with_mobile_and_purpose",
			&IndividualInfo{
				BakongAccountID: "jonhsmith@nbcq",
				MerchantName:    "Jonh Smith",
				MobileNumber:    "85512345678",
				Purpose:         "Testing",
			},
			"00020101021129180014jonhsmith@nbcq5204599953031165802KH5910Jonh Smith6010Phnom Penh62260211855123456780807Testing6304A586",
		},
		{
			"static_with_alt_language_full",
			&IndividualInfo{
				BakongAccountID:       "jonhsmith@nbcq",
				MerchantName:          "Jonh Smith",
				MerchantCity:          "Siam Reap",
				AltLanguagePreference: "km",
				AltMerchantName:       "ចន ស្មីន",
				AltMerchantCity:       "សៀមរាប",
			},
			"00020101021129180014jonhsmith@nbcq5204599953031165802KH5910Jonh Smith6009Siam Reap64280002km0108ចន ស្មីន0206សៀមរាប6304A8B1",
		},
		{
			"static_with_alt_language_name_only",
			&IndividualInfo{
				BakongAccountID:       "jonhsmith@nbcq",
				MerchantName:          "Jonh Smith",
				AltLanguagePreference: "km",
				AltMerchantName:       "ចន ស្មីន",
			},
			"00020101021129180014jonhsmith@nbcq5204599953031165802KH5910Jonh Smith6010Phnom Penh64180002km0108ចន ស្មីន6304A291",
		},
		{
			"static_with_upi",
			&IndividualInfo{
				BakongAccountID: "jonhsmith@nbcq",
				MerchantName:    "Jonh Smith",
				UPIAccountInfo:  "12345678123456789012345",
			},
			"00020101021115231234567812345678901234529180014jonhsmith@nbcq5204599953031165802KH5910Jonh Smith6010Phnom Penh6304F646",
		},
		{
			"static_with_all_account_fields",
			&IndividualInfo{
				BakongAccountID: "jonhsmith@nbcq",
				AccountInfo:     "012345678",
				AcquiringBank:   "Dev Bank",
				MerchantName:    "ចន ស្មីន",
				MerchantCity:    "Siam Reap",
				BillNumber:      "INV-2021-07-65822",
				MobileNumber:    "85512345678",
				StoreLabel:      "BKK-1",
				TerminalLabel:   "012345",
			},
			"00020101021129430014jonhsmith@nbcq01090123456780208Dev Bank5204599953031165802KH5908ចន ស្មីន6009Siam Reap62550117INV-2021-07-658220211855123456780305BKK-107060123456304B956",
		},
		{
			"USD_amount_0.5",
			&IndividualInfo{
				BakongAccountID:     "jonhsmith@nbcq",
				MerchantName:        "Jonh Smith",
				Currency:            USD,
				Amount:              0.5,
				ExpirationTimestamp: futureTS,
			},
			"00020101021229180014jonhsmith@nbcq52045999530384054030.55802KH5910Jonh Smith6010Phnom Penh993400131726749416311011317267500163116304275E",
		},
		{
			"USD_amount_0.55",
			&IndividualInfo{
				BakongAccountID:     "jonhsmith@nbcq",
				MerchantName:        "Jonh Smith",
				Currency:            USD,
				Amount:              0.55,
				ExpirationTimestamp: futureTS,
			},
			"00020101021229180014jonhsmith@nbcq52045999530384054040.555802KH5910Jonh Smith6010Phnom Penh993400131726749416313011317267500163136304A72A",
		},
		{
			"static_max_mobile_number",
			&IndividualInfo{
				BakongAccountID: "jonhsmith@nbcq",
				MerchantName:    "Jonh Smith",
				MobileNumber:    "0123456789012345678901234",
			},
			"00020101021129180014jonhsmith@nbcq5204599953031165802KH5910Jonh Smith6010Phnom Penh622902250123456789012345678901234630402B7",
		},
		{
			"static_category_code_0000",
			&IndividualInfo{
				BakongAccountID:      "jonhsmith@nbcq",
				MerchantName:         "Jonh Smith",
				MerchantCategoryCode: "0000",
			},
			"00020101021129180014jonhsmith@nbcq5204000053031165802KH5910Jonh Smith6010Phnom Penh63045D73",
		},
		{
			"static_category_code_1234",
			&IndividualInfo{
				BakongAccountID:      "jonhsmith@nbcq",
				MerchantName:         "Jonh Smith",
				MerchantCategoryCode: "1234",
			},
			"00020101021129180014jonhsmith@nbcq5204123453031165802KH5910Jonh Smith6010Phnom Penh6304E055",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data, err := GenerateIndividual(tt.info)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.info.Amount > 0 {
				// Dynamic QR: strip timestamp (tag 99) + CRC suffix (46 chars)
				got := data.QR[:len(data.QR)-46]
				want := tt.want[:len(tt.want)-46]
				if got != want {
					t.Errorf("QR prefix mismatch\ngot:  %s\nwant: %s", got, want)
				}
			} else if data.QR != tt.want {
				t.Errorf("QR mismatch\ngot:  %s\nwant: %s", data.QR, tt.want)
			}
		})
	}
}

// --- Generate Merchant Success ---

func TestGenerateMerchantSuccess(t *testing.T) {
	t.Parallel()

	futureTS := time.Now().Add(2 * time.Minute).UnixMilli()

	tests := []struct {
		name string
		info *MerchantInfo
		want string
	}{
		{
			"KHR_all_additional_data",
			&MerchantInfo{
				BakongAccountID:     "jonhsmith@devb",
				MerchantID:          "123456",
				AcquiringBank:       "Dev Bank",
				MerchantName:        "Jonh Smith",
				MerchantCity:        "Siam Reap",
				Currency:            KHR,
				Amount:              50000,
				BillNumber:          "INV-2021-07-65822",
				MobileNumber:        "85512345678",
				StoreLabel:          "BKK-1",
				TerminalLabel:       "012345",
				ExpirationTimestamp: futureTS,
			},
			"00020101021230400014jonhsmith@devb01061234560208Dev Bank5204599953031165405500005802KH5910Jonh Smith6009Siam Reap62550117INV-2021-07-658220211855123456780305BKK-10706012345993400131726749812358011317267504123586304E971",
		},
		{
			"USD_with_store_label",
			&MerchantInfo{
				BakongAccountID:     "jonhsmith@devb",
				MerchantID:          "123456",
				AcquiringBank:       "Dev Bank",
				MerchantName:        "Jonh Smith",
				MerchantCity:        "Phnom Penh",
				Currency:            USD,
				Amount:              10,
				BillNumber:          "INV-2021-07-65822",
				MobileNumber:        "85512345678",
				StoreLabel:          "BKK-1",
				ExpirationTimestamp: futureTS,
			},
			"00020101021230400014jonhsmith@devb01061234560208Dev Bank5204599953038405402105802KH5910Jonh Smith6010Phnom Penh62450117INV-2021-07-658220211855123456780305BKK-19934001317267498124020113172675041240263047D38",
		},
		{
			"USD_bill_and_mobile",
			&MerchantInfo{
				BakongAccountID:     "jonhsmith@devb",
				MerchantID:          "123456",
				AcquiringBank:       "Dev Bank",
				MerchantName:        "Jonh Smith",
				MerchantCity:        "Phnom Penh",
				Currency:            USD,
				Amount:              10,
				BillNumber:          "INV-2021-07-65822",
				MobileNumber:        "85512345678",
				ExpirationTimestamp: futureTS,
			},
			"00020101021230400014jonhsmith@devb01061234560208Dev Bank5204599953038405402105802KH5910Jonh Smith6010Phnom Penh62360117INV-2021-07-658220211855123456789934001317267498124070113172675041240763040441",
		},
		{
			"USD_bill_only",
			&MerchantInfo{
				BakongAccountID:     "jonhsmith@devb",
				MerchantID:          "123456",
				AcquiringBank:       "Dev Bank",
				MerchantName:        "Jonh Smith",
				MerchantCity:        "Phnom Penh",
				Currency:            USD,
				Amount:              10,
				BillNumber:          "INV-2021-07-65822",
				ExpirationTimestamp: futureTS,
			},
			"00020101021230400014jonhsmith@devb01061234560208Dev Bank5204599953038405402105802KH5910Jonh Smith6010Phnom Penh62210117INV-2021-07-65822993400131726749812411011317267504124116304506A",
		},
		{
			"USD_no_additional_data",
			&MerchantInfo{
				BakongAccountID:     "jonhsmith@devb",
				MerchantID:          "123456",
				AcquiringBank:       "Dev Bank",
				MerchantName:        "Jonh Smith",
				MerchantCity:        "Phnom Penh",
				Currency:            USD,
				Amount:              10,
				ExpirationTimestamp: futureTS,
			},
			"00020101021230400014jonhsmith@devb01061234560208Dev Bank5204599953038405402105802KH5910Jonh Smith6010Phnom Penh99340013172674981241501131726750412415630429F4",
		},
		{
			"static_USD_no_amount",
			&MerchantInfo{
				BakongAccountID: "jonhsmith@devb",
				MerchantID:      "123456",
				AcquiringBank:   "Dev Bank",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "Phnom Penh",
				Currency:        USD,
			},
			"00020101021130400014jonhsmith@devb01061234560208Dev Bank5204599953038405802KH5910Jonh Smith6010Phnom Penh630484CA",
		},
		{
			"USD_khmer_name",
			&MerchantInfo{
				BakongAccountID:     "jonhsmith@devb",
				MerchantID:          "123456",
				AcquiringBank:       "Dev Bank",
				MerchantName:        "ចន ស្មីន",
				MerchantCity:        "Phnom Penh",
				Currency:            USD,
				Amount:              1.1,
				ExpirationTimestamp: futureTS,
			},
			"00020101021230400014jonhsmith@devb01061234560208Dev Bank52045999530384054031.15802KH5908ចន ស្មីន6010Phnom Penh993400131726749812422011317267504124226304D63E",
		},
		{
			"static_with_purpose",
			&MerchantInfo{
				BakongAccountID: "jonhsmith@nbcq",
				MerchantID:      "123456",
				AcquiringBank:   "Dev Bank",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "Phnom Penh",
				Purpose:         "Testing",
			},
			"00020101021130400014jonhsmith@nbcq01061234560208Dev Bank5204599953031165802KH5910Jonh Smith6010Phnom Penh62110807Testing63041643",
		},
		{
			"static_with_alt_language_full",
			&MerchantInfo{
				BakongAccountID:       "jonhsmith@nbcq",
				MerchantID:            "123456",
				AcquiringBank:         "Dev Bank",
				MerchantName:          "Jonh Smith",
				MerchantCity:          "Siam Reap",
				AltLanguagePreference: "km",
				AltMerchantName:       "ចន ស្មីន",
				AltMerchantCity:       "សៀមរាប",
			},
			"00020101021130400014jonhsmith@nbcq01061234560208Dev Bank5204599953031165802KH5910Jonh Smith6009Siam Reap64280002km0108ចន ស្មីន0206សៀមរាប6304C9AA",
		},
		{
			"static_with_alt_language_name_only",
			&MerchantInfo{
				BakongAccountID:       "jonhsmith@nbcq",
				MerchantID:            "123456",
				AcquiringBank:         "Dev Bank",
				MerchantName:          "Jonh Smith",
				MerchantCity:          "Phnom Penh",
				AltLanguagePreference: "km",
				AltMerchantName:       "ចន ស្មីន",
			},
			"00020101021130400014jonhsmith@nbcq01061234560208Dev Bank5204599953031165802KH5910Jonh Smith6010Phnom Penh64180002km0108ចន ស្មីន6304A4D1",
		},
		{
			"static_with_upi",
			&MerchantInfo{
				BakongAccountID: "jonhsmith@nbcq",
				MerchantID:      "123456",
				AcquiringBank:   "Dev Bank",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "Phnom Penh",
				UPIAccountInfo:  "12345678123456789012345",
			},
			"00020101021115231234567812345678901234530400014jonhsmith@nbcq01061234560208Dev Bank5204599953031165802KH5910Jonh Smith6010Phnom Penh63049FF0",
		},
		{
			"static_category_code_0000",
			&MerchantInfo{
				BakongAccountID:      "jonhsmith@nbcq",
				MerchantID:           "123456",
				AcquiringBank:        "Dev Bank",
				MerchantName:         "Jonh Smith",
				MerchantCity:         "Phnom Penh",
				MerchantCategoryCode: "0000",
			},
			"00020101021130400014jonhsmith@nbcq01061234560208Dev Bank5204000053031165802KH5910Jonh Smith6010Phnom Penh630483B1",
		},
		{
			"static_category_code_1234",
			&MerchantInfo{
				BakongAccountID:      "jonhsmith@nbcq",
				MerchantID:           "123456",
				AcquiringBank:        "Dev Bank",
				MerchantName:         "Jonh Smith",
				MerchantCity:         "Phnom Penh",
				MerchantCategoryCode: "1234",
			},
			"00020101021130400014jonhsmith@nbcq01061234560208Dev Bank5204123453031165802KH5910Jonh Smith6010Phnom Penh63043E97",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data, err := GenerateMerchant(tt.info)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.info.Amount > 0 {
				got := data.QR[:len(data.QR)-46]
				want := tt.want[:len(tt.want)-46]
				if got != want {
					t.Errorf("QR prefix mismatch\ngot:  %s\nwant: %s", got, want)
				}
			} else if data.QR != tt.want {
				t.Errorf("QR mismatch\ngot:  %s\nwant: %s", data.QR, tt.want)
			}
		})
	}
}

// --- Generate Individual Required Errors ---

func TestGenerateIndividualRequiredError(t *testing.T) {
	t.Parallel()

	futureTS := time.Now().Add(2 * time.Minute).UnixMilli()
	_ = futureTS

	tests := []struct {
		name    string
		info    *IndividualInfo
		wantErr error
	}{
		{
			"missing_account_id",
			&IndividualInfo{MerchantName: "Jonh Smith"},
			ErrAccountIDRequired,
		},
		{
			"missing_merchant_name",
			&IndividualInfo{BakongAccountID: "johnsmith@devb"},
			ErrMerchantNameRequired,
		},
		{
			"missing_language_preference_with_alt_name_and_city",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				AltMerchantName: "ចន ស្មីន",
				AltMerchantCity: "សៀមរាប",
			},
			ErrLanguagePreferenceRequired,
		},
		{
			"missing_alt_merchant_name",
			&IndividualInfo{
				BakongAccountID:       "johnsmith@devb",
				MerchantName:          "Jonh Smith",
				AltLanguagePreference: "km",
				AltMerchantCity:       "សៀមរាប",
			},
			ErrMerchantNameAltRequired,
		},
		{
			"missing_language_preference_with_alt_city_only",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				AltMerchantCity: "សៀមរាប",
			},
			ErrLanguagePreferenceRequired,
		},
		{
			"missing_expiration_with_amount",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				Currency:        KHR,
				Amount:          101,
			},
			ErrExpirationRequired,
		},
		{
			"expiration_in_past",
			&IndividualInfo{
				BakongAccountID:     "johnsmith@devb",
				MerchantName:        "Jonh Smith",
				Currency:            USD,
				Amount:              101.3,
				ExpirationTimestamp: 1727260807000,
			},
			ErrExpirationInPast,
		},
		{
			"expiration_invalid_length",
			&IndividualInfo{
				BakongAccountID:     "johnsmith@devb",
				MerchantName:        "Jonh Smith",
				Currency:            USD,
				Amount:              101.3,
				ExpirationTimestamp: 1727260807,
			},
			ErrInvalidTimestamp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := GenerateIndividual(tt.info)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

// --- Generate Individual Length/Validation Errors ---

func TestGenerateIndividualLengthError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		info    *IndividualInfo
		wantErr error
	}{
		{
			"account_id_too_long",
			&IndividualInfo{
				BakongAccountID: "@johnsmith00123456789012345678912345@devb",
				MerchantName:    "Jonh Smith",
				Currency:        USD,
			},
			ErrAccountIDTooLong,
		},
		{
			"merchant_name_too_long",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "012345678901234567890123456789",
				Currency:        USD,
			},
			ErrMerchantNameTooLong,
		},
		{
			"amount_too_large",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				Currency:        USD,
				Amount:          123456789012345,
			},
			ErrInvalidAmount,
		},
		{
			"merchant_city_too_long",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "012345678901234567890123456789",
				Currency:        USD,
				Amount:          123,
			},
			ErrMerchantCityTooLong,
		},
		{
			"bill_number_too_long",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "PP",
				Currency:        USD,
				Amount:          100,
				BillNumber:      "012345678901234567890123456789",
			},
			ErrBillNumberTooLong,
		},
		{
			"mobile_number_too_long",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "PP",
				Currency:        USD,
				Amount:          100,
				MobileNumber:    "12345678901234567890123456",
			},
			ErrMobileNumberTooLong,
		},
		{
			"store_label_too_long",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "PP",
				Currency:        USD,
				Amount:          100,
				StoreLabel:      "012345678901234567890123456789",
			},
			ErrStoreLabelTooLong,
		},
		{
			"terminal_label_too_long",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "PP",
				Currency:        USD,
				Amount:          100,
				TerminalLabel:   "012345678901234567890123456789",
			},
			ErrTerminalLabelTooLong,
		},
		{
			"acquiring_bank_too_long",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				AcquiringBank:   "Advanced Bank of Asia Limited Cambodia",
				MerchantName:    "ABC",
				Currency:        USD,
				Amount:          100,
				MerchantCity:    "Phnom Penh",
			},
			ErrAcquiringBankTooLong,
		},
		{
			"account_info_too_long",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				AccountInfo:     "012345678901234567890123456789897",
				MerchantName:    "ABC",
				Currency:        USD,
				Amount:          100,
			},
			ErrAccountInfoTooLong,
		},
		{
			"usd_three_decimal_places",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "ABC",
				Currency:        USD,
				Amount:          1.234,
			},
			ErrInvalidAmount,
		},
		{
			"upi_too_long",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "ABC",
				Currency:        KHR,
				Amount:          10,
				UPIAccountInfo:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis a ultrices nunc. Quisque in sem fringilla, ullamcorper est vel",
			},
			ErrUPITooLong,
		},
		{
			"language_preference_too_long",
			&IndividualInfo{
				BakongAccountID:       "johnsmith@devb",
				MerchantName:          "ABC",
				Currency:              USD,
				Amount:                10,
				AltLanguagePreference: "khmer",
			},
			ErrLanguagePreferenceTooLong,
		},
		{
			"merchant_name_alt_too_long",
			&IndividualInfo{
				BakongAccountID:       "johnsmith@devb",
				MerchantName:          "ABC",
				Currency:              USD,
				Amount:                10,
				AltLanguagePreference: "km",
				AltMerchantName:       "Lorem ipsum dolor sit amet",
			},
			ErrMerchantNameAltTooLong,
		},
		{
			"merchant_city_alt_too_long",
			&IndividualInfo{
				BakongAccountID:       "johnsmith@devb",
				MerchantName:          "ABC",
				Currency:              USD,
				Amount:                10,
				AltLanguagePreference: "km",
				AltMerchantName:       "Lorem ipsum dolor sit am",
				AltMerchantCity:       "Lorem ipsum dolor",
			},
			ErrMerchantCityAltTooLong,
		},
		{
			"purpose_too_long",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "ABC",
				Currency:        USD,
				Amount:          10,
				Purpose:         "Lorem ipsum dolor sit amet",
			},
			ErrPurposeTooLong,
		},
		{
			"invalid_merchant_category_code_non_numeric",
			&IndividualInfo{
				BakongAccountID:      "johnsmith@devb",
				MerchantName:         "Jonh Smith",
				Currency:             USD,
				MerchantCategoryCode: "1A2B",
			},
			ErrMerchantCategoryCodeInvalid,
		},
		{
			"upi_not_supported_with_usd",
			&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "ABC",
				Currency:        USD,
				UPIAccountInfo:  "012345678901234567890123456789",
			},
			ErrUPINotSupportUSD,
		},
		{
			"invalid_merchant_category_code_negative",
			&IndividualInfo{
				BakongAccountID:      "johnsmith@devb",
				MerchantName:         "Jonh Smith",
				Currency:             USD,
				MerchantCategoryCode: "-100",
			},
			ErrMerchantCategoryCodeInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := GenerateIndividual(tt.info)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

// --- Generate Merchant Required Errors ---

func TestGenerateMerchantRequiredError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		info    *MerchantInfo
		wantErr error
	}{
		{
			"missing_account_id",
			&MerchantInfo{},
			ErrAccountIDRequired,
		},
		{
			"missing_merchant_id",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
			},
			ErrMerchantIDRequired,
		},
		{
			"missing_acquiring_bank",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				MerchantID:      "12345",
			},
			ErrAcquiringBankRequired,
		},
		{
			"missing_merchant_name",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantID:      "12345",
				AcquiringBank:   "ABA",
				MerchantCity:    "Phnom Penh",
			},
			ErrMerchantNameRequired,
		},
		{
			"missing_merchant_id_with_bank",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				AcquiringBank:   "Dev Bank",
			},
			ErrMerchantIDRequired,
		},
		{
			"missing_acquiring_bank_with_id",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				MerchantID:      "123456",
			},
			ErrAcquiringBankRequired,
		},
		{
			"missing_language_preference_with_alt_name_and_city",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				MerchantID:      "123456",
				AcquiringBank:   "Dev Bank",
				MerchantCity:    "Phnom Penh",
				AltMerchantName: "ចន ស្មីន",
				AltMerchantCity: "សៀមរាប",
			},
			ErrLanguagePreferenceRequired,
		},
		{
			"missing_language_preference_with_alt_city_only",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				MerchantID:      "123456",
				AcquiringBank:   "Dev Bank",
				MerchantCity:    "Phnom Penh",
				AltMerchantCity: "សៀមរាប",
			},
			ErrLanguagePreferenceRequired,
		},
		{
			"missing_alt_merchant_name",
			&MerchantInfo{
				BakongAccountID:       "johnsmith@devb",
				MerchantName:          "Jonh Smith",
				MerchantID:            "123456",
				AcquiringBank:         "Dev Bank",
				MerchantCity:          "Phnom Penh",
				AltLanguagePreference: "km",
				AltMerchantCity:       "សៀមរាប",
			},
			ErrMerchantNameAltRequired,
		},
		{
			"missing_expiration_with_amount",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				MerchantID:      "123456",
				AcquiringBank:   "Dev Bank",
				MerchantCity:    "Phnom Penh",
				Currency:        USD,
				Amount:          102.3,
			},
			ErrExpirationRequired,
		},
		{
			"expiration_in_past",
			&MerchantInfo{
				BakongAccountID:     "johnsmith@devb",
				MerchantName:        "Jonh Smith",
				MerchantID:          "123456",
				AcquiringBank:       "Dev Bank",
				MerchantCity:        "Phnom Penh",
				Currency:            USD,
				Amount:              102.3,
				ExpirationTimestamp: 1727260807000,
			},
			ErrExpirationInPast,
		},
		{
			"expiration_invalid_length",
			&MerchantInfo{
				BakongAccountID:     "johnsmith@devb",
				MerchantName:        "Jonh Smith",
				MerchantID:          "123456",
				AcquiringBank:       "Dev Bank",
				MerchantCity:        "Phnom Penh",
				Currency:            USD,
				Amount:              102.3,
				ExpirationTimestamp: 1727260807,
			},
			ErrInvalidTimestamp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := GenerateMerchant(tt.info)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

// --- Generate Merchant Length/Validation Errors ---

func TestGenerateMerchantLengthError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		info    *MerchantInfo
		wantErr error
	}{
		{
			"account_id_too_long",
			&MerchantInfo{
				BakongAccountID: "johnsmith00123456789012345678912345@devb",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "Phnom Penh",
				Currency:        USD,
			},
			ErrAccountIDTooLong,
		},
		{
			"merchant_id_too_long",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantID:      "ohnsmith001234567890123456789123450123456",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "Phnom Penh",
				Currency:        USD,
			},
			ErrMerchantIDTooLong,
		},
		{
			"acquiring_bank_too_long",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantID:      "123",
				AcquiringBank:   "ohnsmith001234567890123456789123450123456",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "Phnom Penh",
				Currency:        USD,
			},
			ErrAcquiringBankTooLong,
		},
		{
			"merchant_name_too_long",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantID:      "123",
				AcquiringBank:   "Dev Bank",
				MerchantName:    "ohnsmith001234567890123456789123450123456",
				MerchantCity:    "Phnom Penh",
				Currency:        USD,
			},
			ErrMerchantNameTooLong,
		},
		{
			"amount_too_large",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantID:      "123",
				AcquiringBank:   "Dev Bank",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "Phnom Penh",
				Currency:        USD,
				Amount:          123456789012345,
			},
			ErrInvalidAmount,
		},
		{
			"merchant_city_too_long",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantID:      "123",
				AcquiringBank:   "Dev Bank",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "ohnsmith001234567890123456789123450123456",
				Currency:        USD,
				Amount:          100,
			},
			ErrMerchantCityTooLong,
		},
		{
			"bill_number_too_long",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantID:      "123",
				AcquiringBank:   "Dev Bank",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "Phnom penh",
				Currency:        USD,
				Amount:          100,
				BillNumber:      "ohnsmith001234567890123456789123450123456",
			},
			ErrBillNumberTooLong,
		},
		{
			"mobile_number_too_long",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantID:      "123",
				AcquiringBank:   "Dev Bank",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "Phnom penh",
				Currency:        USD,
				Amount:          100,
				BillNumber:      "123",
				MobileNumber:    "12345678901234567890123456",
			},
			ErrMobileNumberTooLong,
		},
		{
			"store_label_too_long",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantID:      "123",
				AcquiringBank:   "Dev Bank",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "Phnom penh",
				Currency:        USD,
				Amount:          100,
				BillNumber:      "123",
				MobileNumber:    "855123172",
				StoreLabel:      "ohnsmith001234567890123456789123450123456",
			},
			ErrStoreLabelTooLong,
		},
		{
			"terminal_label_too_long",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantID:      "123",
				AcquiringBank:   "Dev Bank",
				MerchantName:    "Jonh Smith",
				MerchantCity:    "Phnom penh",
				Currency:        USD,
				Amount:          100,
				BillNumber:      "123",
				MobileNumber:    "855123172",
				StoreLabel:      "123",
				TerminalLabel:   "ohnsmith001234567890123456789123450123456",
			},
			ErrTerminalLabelTooLong,
		},
		{
			"upi_too_long",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantID:      "123456",
				AcquiringBank:   "Dev Bank",
				MerchantName:    "ABC",
				MerchantCity:    "Phnom Penh",
				Currency:        KHR,
				Amount:          10,
				UPIAccountInfo:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis a ultrices nunc. Quisque in sem fringilla, ullamcorper est vel",
			},
			ErrUPITooLong,
		},
		{
			"language_preference_too_long",
			&MerchantInfo{
				BakongAccountID:       "johnsmith@devb",
				MerchantID:            "123456",
				AcquiringBank:         "Dev Bank",
				MerchantName:          "ABC",
				MerchantCity:          "Phnom Penh",
				Currency:              USD,
				Amount:                10,
				AltLanguagePreference: "Khmer",
			},
			ErrLanguagePreferenceTooLong,
		},
		{
			"merchant_name_alt_too_long",
			&MerchantInfo{
				BakongAccountID:       "johnsmith@devb",
				MerchantID:            "123456",
				AcquiringBank:         "Dev Bank",
				MerchantName:          "ABC",
				MerchantCity:          "Phnom Penh",
				Currency:              USD,
				Amount:                10,
				AltLanguagePreference: "km",
				AltMerchantName:       "Lorem ipsum dolor sit amet",
			},
			ErrMerchantNameAltTooLong,
		},
		{
			"merchant_city_alt_too_long",
			&MerchantInfo{
				BakongAccountID:       "johnsmith@devb",
				MerchantID:            "123456",
				AcquiringBank:         "Dev Bank",
				MerchantName:          "ABC",
				MerchantCity:          "Phnom Penh",
				Currency:              USD,
				Amount:                10,
				AltLanguagePreference: "km",
				AltMerchantName:       "Lorem ipsum dolor",
				AltMerchantCity:       "Lorem ipsum dolor",
			},
			ErrMerchantCityAltTooLong,
		},
		{
			"purpose_too_long",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantID:      "123456",
				AcquiringBank:   "Dev Bank",
				MerchantName:    "ABC",
				MerchantCity:    "Phnom Penh",
				Currency:        USD,
				Amount:          10,
				Purpose:         "Lorem ipsum dolor sit amet",
			},
			ErrPurposeTooLong,
		},
		{
			"invalid_merchant_category_code_non_numeric",
			&MerchantInfo{
				BakongAccountID:      "johnsmith@devb",
				MerchantID:           "123456",
				AcquiringBank:        "Dev Bank",
				MerchantName:         "Jonh Smith",
				MerchantCity:         "Phnom Penh",
				Currency:             USD,
				MerchantCategoryCode: "1A2B",
			},
			ErrMerchantCategoryCodeInvalid,
		},
		{
			"upi_not_supported_with_usd",
			&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantID:      "123456",
				AcquiringBank:   "Dev Bank",
				MerchantName:    "ABC",
				MerchantCity:    "Phnom Penh",
				Currency:        USD,
				Amount:          10,
				UPIAccountInfo:  "12345678901234567890123456",
			},
			ErrUPINotSupportUSD,
		},
		{
			"invalid_merchant_category_code_negative",
			&MerchantInfo{
				BakongAccountID:      "johnsmith@devb",
				MerchantID:           "123456",
				AcquiringBank:        "Dev Bank",
				MerchantName:         "Jonh Smith",
				MerchantCity:         "Phnom Penh",
				Currency:             USD,
				MerchantCategoryCode: "-100",
			},
			ErrMerchantCategoryCodeInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := GenerateMerchant(tt.info)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

// --- Amount Validation: Invalid Amounts ---

func TestInvalidAmountIndividual(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		currency Currency
		amount   float64
	}{
		{"USD_three_decimals", USD, 1.234},
		{"USD_negative", USD, -1000},
		{"USD_five_decimals", USD, 100.00111},
		{"USD_amount_too_large", USD, 12345678901234},
		{"USD_amount_exceeds_max", USD, 999999999999.99},
		{"KHR_negative", KHR, -1000},
		{"KHR_five_decimals", KHR, 100.00111},
		{"KHR_amount_too_large", KHR, 12345678901234},
		{"KHR_decimal_999999999999.99", KHR, 999999999999.99},
		{"KHR_one_decimal", KHR, 1.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := GenerateIndividual(&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				Currency:        tt.currency,
				Amount:          tt.amount,
			})
			if !errors.Is(err, ErrInvalidAmount) {
				t.Errorf("expected %v, got %v", ErrInvalidAmount, err)
			}
		})
	}
}

func TestInvalidAmountMerchant(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		currency Currency
		amount   float64
	}{
		{"USD_three_decimals", USD, 1.234},
		{"USD_negative", USD, -1000},
		{"USD_five_decimals", USD, 100.00111},
		{"USD_amount_exceeds_max", USD, 999999999999.99},
		{"KHR_negative", KHR, -1000},
		{"KHR_five_decimals", KHR, 100.00111},
		{"KHR_amount_too_large", KHR, 12345678901234},
		{"KHR_decimal_999999999999.99", KHR, 999999999999.99},
		{"KHR_one_decimal", KHR, 1.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := GenerateMerchant(&MerchantInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				MerchantID:      "123456",
				AcquiringBank:   "Dev Bank",
				MerchantCity:    "PP",
				Currency:        tt.currency,
				Amount:          tt.amount,
			})
			if !errors.Is(err, ErrInvalidAmount) {
				t.Errorf("expected %v, got %v", ErrInvalidAmount, err)
			}
		})
	}
}

// --- Amount Formatting: Exact Amount Values ---

func TestExactAmountIndividual(t *testing.T) {
	t.Parallel()

	futureTS := time.Now().Add(2 * time.Minute).UnixMilli()

	tests := []struct {
		name       string
		currency   Currency
		amount     float64
		wantAmount string
	}{
		{"KHR_100", KHR, 100, "100"},
		{"KHR_9999999999", KHR, 9999999999, "9999999999"},
		{"KHR_10000", KHR, 10000, "10000"},
		{"KHR_1234567890123", KHR, 1234567890123, "1234567890123"},
		{"USD_1.12", USD, 1.12, "1.12"},
		{"USD_1000", USD, 1000, "1000"},
		{"USD_100.11", USD, 100.11, "100.11"},
		{"USD_100.12_trailing_zeros", USD, 100.12, "100.12"},
		{"USD_12345678901", USD, 12345678901, "12345678901"},
		{"USD_9999999999.99", USD, 9999999999.99, "9999999999.99"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data, err := GenerateIndividual(&IndividualInfo{
				BakongAccountID:     "johnsmith@devb",
				AccountInfo:         "012345678",
				AcquiringBank:       "Dev Bank",
				MerchantName:        "Jonh Smith",
				MerchantCity:        "Siam Reap",
				Currency:            tt.currency,
				Amount:              tt.amount,
				ExpirationTimestamp: futureTS,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			decoded, err := Decode(data.QR)
			if err != nil {
				t.Fatalf("decode error: %v", err)
			}
			if decoded.TransactionAmount != tt.wantAmount {
				t.Errorf("amount = %q, want %q", decoded.TransactionAmount, tt.wantAmount)
			}
		})
	}
}

func TestExactAmountMerchant(t *testing.T) {
	t.Parallel()

	futureTS := time.Now().Add(2 * time.Minute).UnixMilli()

	tests := []struct {
		name       string
		currency   Currency
		amount     float64
		wantAmount string
	}{
		{"USD_1.1", USD, 1.1, "1.1"},
		{"USD_1.12", USD, 1.12, "1.12"},
		{"USD_100.12_trailing_zeros", USD, 100.12, "100.12"},
		{"USD_12345678901", USD, 12345678901, "12345678901"},
		{"USD_9999999999.99", USD, 9999999999.99, "9999999999.99"},
		{"KHR_100", KHR, 100, "100"},
		{"KHR_9999999999", KHR, 9999999999, "9999999999"},
		{"KHR_10000", KHR, 10000, "10000"},
		{"KHR_1234567890123", KHR, 1234567890123, "1234567890123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data, err := GenerateMerchant(&MerchantInfo{
				BakongAccountID:     "johnsmith@devb",
				MerchantID:          "0123456",
				AcquiringBank:       "Dev Bank",
				MerchantName:        "Jonh Smith",
				MerchantCity:        "Siam Reap",
				Currency:            tt.currency,
				Amount:              tt.amount,
				ExpirationTimestamp: futureTS,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			decoded, err := Decode(data.QR)
			if err != nil {
				t.Fatalf("decode error: %v", err)
			}
			if decoded.TransactionAmount != tt.wantAmount {
				t.Errorf("amount = %q, want %q", decoded.TransactionAmount, tt.wantAmount)
			}
		})
	}
}

// --- Amount Validation: Valid Amounts ---

func TestValidAmount(t *testing.T) {
	t.Parallel()

	futureTS := time.Now().Add(2 * time.Minute).UnixMilli()

	tests := []struct {
		name     string
		amount   float64
		currency Currency
	}{
		{"KHR_1000", 1000, KHR},
		{"KHR_100", 100, KHR},
		{"KHR_9999999999", 9999999999, KHR},
		{"KHR_1", 1, KHR},
		{"USD_1000", 1000, USD},
		{"USD_100.11", 100.11, USD},
		{"USD_100.12", 100.12, USD},
		{"USD_100.01", 100.01, USD},
		{"USD_12345678901", 12345678901, USD},
		{"USD_9999999999.99", 9999999999.99, USD},
		{"USD_1.12", 1.12, USD},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := GenerateIndividual(&IndividualInfo{
				BakongAccountID:     "johnsmith@devb",
				MerchantName:        "Jonh Smith",
				Currency:            tt.currency,
				Amount:              tt.amount,
				ExpirationTimestamp: futureTS,
			})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// --- Tag Validation: Invalid Bakong Account IDs ---

func TestInvalidBakongAccountID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		accountID string
	}{
		{"hash_instead_of_at", "abcd#abcd"},
		{"no_at_sign", "abcdefgh"},
		{"at_prefix", "@abcd"},
		{"at_suffix", "abcd@"},
		{"double_at_suffix", "abcd@@"},
		{"double_at_prefix", "@@abcd"},
		{"all_at_signs", "@@@@"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := GenerateIndividual(&IndividualInfo{
				BakongAccountID: tt.accountID,
				MerchantName:    "Jonh Smith",
			})
			if !errors.Is(err, ErrAccountIDInvalid) {
				t.Errorf("expected %v, got %v", ErrAccountIDInvalid, err)
			}
		})
	}
}

// --- Tag Validation: Additional Data Tag Assignment ---

func TestAdditionalDataTagAssignment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		billNumber    string
		terminalLabel string
		storeLabel    string
		wrongSnippet  string // snippet with swapped subtag assignments
	}{
		{
			"terminal_not_encoded_as_store",
			"#12345",
			"BKK Store",
			"",
			"62230106#123450309BKK Store",
		},
		{
			"store_and_terminal_not_swapped",
			"#12345",
			"BKK Store",
			"#2",
			"62290106#123450309BKK Store0702#2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data, err := GenerateIndividual(&IndividualInfo{
				BakongAccountID: "johnsmith@devb",
				MerchantName:    "Jonh Smith",
				BillNumber:      tt.billNumber,
				TerminalLabel:   tt.terminalLabel,
				StoreLabel:      tt.storeLabel,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if strings.Contains(data.QR, tt.wrongSnippet) {
				t.Errorf("QR should not contain wrong tag assignment %q\nQR: %s", tt.wrongSnippet, data.QR)
			}
		})
	}
}
