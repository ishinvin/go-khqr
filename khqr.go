// Package khqr provides functionality to generate, decode, and verify
// KHQR payment codes used with Cambodia's Bakong payment system.
package khqr

// GenerateIndividual generates a KHQR string for an individual payment.
func GenerateIndividual(info IndividualInfo) (*Data, error) { //nolint:gocritic // value param creates shallow copy to avoid mutating caller's struct
	return generateIndividual(&info)
}

// GenerateMerchant generates a KHQR string for a merchant payment.
func GenerateMerchant(info MerchantInfo) (*Data, error) { //nolint:gocritic // value param creates shallow copy to avoid mutating caller's struct
	return generateMerchant(&info)
}

// Decode parses a KHQR string and returns structured data.
func Decode(qr string) (*DecodedData, error) {
	return decode(qr)
}

// Verify validates the CRC and structure of a KHQR string.
// Returns nil if valid, or an error describing the issue.
func Verify(qr string) error {
	return verify(qr)
}
