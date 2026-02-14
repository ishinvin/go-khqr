package khqr

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// tlv represents a single Tag-Length-Value entry.
type tlv struct {
	Tag   string
	Value string
}

// tlvWriter wraps strings.Builder with a helper for conditional TLV encoding.
type tlvWriter struct {
	strings.Builder
}

// writeTLV encodes a tag-value pair and writes it to the builder.
// If value is empty, nothing is written.
func (w *tlvWriter) writeTLV(tag, value string) {
	if value != "" {
		w.WriteString(encodeTLV(tag, value))
	}
}

// encodeTLV encodes a tag and value into a TLV string.
// Format: tag (2 chars) + length (2 chars, zero-padded) + value
// Length is the number of Unicode characters (runes) in the value, not bytes.
func encodeTLV(tag, value string) string {
	n := utf8.RuneCountInString(value)
	if n < 10 { //nolint:mnd // zero-pad single digit
		return tag + "0" + strconv.Itoa(n) + value
	}
	return tag + strconv.Itoa(n) + value
}

// parseTLV parses a TLV-encoded string into an ordered list of entries.
// The length field is interpreted as a rune count (Unicode characters).
func parseTLV(data string) ([]tlv, error) {
	var entries []tlv
	runes := []rune(data)
	pos := 0

	for pos < len(runes) {
		if pos+4 > len(runes) {
			return nil, ErrInvalidQR
		}

		tag := string(runes[pos : pos+2])
		pos += 2

		lengthStr := string(runes[pos : pos+2])
		pos += 2

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, ErrInvalidQR
		}

		if pos+length > len(runes) {
			return nil, ErrInvalidQR
		}

		value := string(runes[pos : pos+length])
		pos += length

		entries = append(entries, tlv{Tag: tag, Value: value})
	}

	return entries, nil
}
