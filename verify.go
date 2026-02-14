package khqr

import (
	"regexp"
	"strings"
)

var crcFormatRegex = regexp.MustCompile(`6304[A-Fa-f0-9]{4}$`)

// verifyCRC validates the CRC format and checksum of a KHQR string.
func verifyCRC(qr string) error {
	if !crcFormatRegex.MatchString(qr) {
		return ErrCRCInvalid
	}

	crcValue := qr[len(qr)-4:]
	dataForCRC := qr[:len(qr)-4]
	if !strings.EqualFold(crc16Hex(dataForCRC), crcValue) {
		return ErrCRCInvalid
	}

	return nil
}

// verify validates a KHQR string by checking CRC, decoding, and validating all fields.
func verify(qr string) error {
	qr = strings.TrimSpace(qr)
	if len(qr) < 8 { //nolint:mnd // minimum QR length
		return ErrInvalidQR
	}

	if err := verifyCRC(qr); err != nil {
		return err
	}

	data, err := decode(qr)
	if err != nil {
		return err
	}
	return data.validate()
}
