package khqr

import (
	"errors"
	"testing"
)

func TestEncodeTLV(t *testing.T) {
	tests := []struct {
		tag, value, expected string
	}{
		{"00", "01", "000201"},
		{"59", "John Smith", "5910John Smith"},
		{"58", "KH", "5802KH"},
	}

	for _, tt := range tests {
		result := encodeTLV(tt.tag, tt.value)
		if result != tt.expected {
			t.Errorf("encodeTLV(%q, %q) = %q, expected %q", tt.tag, tt.value, result, tt.expected)
		}
	}
}

func TestParseTLV(t *testing.T) {
	entries, err := parseTLV("000201")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Tag != "00" || entries[0].Value != "01" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}

	entries, err = parseTLV("000201" + "010212" + "5802KH")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestParseTLVInvalid(t *testing.T) {
	_, err := parseTLV("00")
	if !errors.Is(err, ErrInvalidQR) {
		t.Errorf("expected ErrInvalidQR, got %v", err)
	}

	_, err = parseTLV("00XX01")
	if !errors.Is(err, ErrInvalidQR) {
		t.Errorf("expected ErrInvalidQR, got %v", err)
	}
}
