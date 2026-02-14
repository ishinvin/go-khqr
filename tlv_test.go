package khqr

import (
	"errors"
	"testing"
)

func TestEncodeTLV(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		tag   string
		value string
		want  string
	}{
		{"empty_value", "00", "", "0000"},
		{"single_char", "01", "A", "0101A"},
		{"single_digit_length", "02", "hi", "0202hi"},
		{"double_digit_length", "03", "0123456789", "03100123456789"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := encodeTLV(tt.tag, tt.value)
			if got != tt.want {
				t.Errorf("encodeTLV(%q, %q) = %q, want %q", tt.tag, tt.value, got, tt.want)
			}
		})
	}
}

func TestParseTLV(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data string
		want []tlv
	}{
		{
			"single_entry",
			"0103abc",
			[]tlv{{Tag: "01", Value: "abc"}},
		},
		{
			"multiple_entries",
			"0103abc" + "0202hi" + "0301X",
			[]tlv{
				{Tag: "01", Value: "abc"},
				{Tag: "02", Value: "hi"},
				{Tag: "03", Value: "X"},
			},
		},
		{
			"empty_value",
			"0100",
			[]tlv{{Tag: "01", Value: ""}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseTLV(tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d entries, want %d", len(got), len(tt.want))
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("entry[%d]: got %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParseTLVInvalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data string
	}{
		{"truncated_header", "01"},
		{"non_numeric_length", "01XXa"},
		{"length_exceeds_data", "0105ab"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := parseTLV(tt.data)
			if !errors.Is(err, ErrInvalidQR) {
				t.Errorf("parseTLV(%q) = %v, want ErrInvalidQR", tt.data, err)
			}
		})
	}
}

func TestTLVWriterWriteTLV(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		tag   string
		value string
		want  string
	}{
		{"non_empty_value", "01", "abc", "0103abc"},
		{"empty_value", "01", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var w tlvWriter
			w.writeTLV(tt.tag, tt.value)
			if got := w.String(); got != tt.want {
				t.Errorf("writeTLV(%q, %q) = %q, want %q", tt.tag, tt.value, got, tt.want)
			}
		})
	}
}

func TestTLVWriterWriteTLVMultiple(t *testing.T) {
	t.Parallel()

	var w tlvWriter
	w.writeTLV("01", "a")
	w.writeTLV("02", "")
	w.writeTLV("03", "bc")

	want := "0101a" + "0302bc"
	if got := w.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
