package khqr

import (
	"errors"
	"testing"
)

func Test_encodeTLV(t *testing.T) {
	tests := []struct {
		name  string
		tag   string
		value string
		want  string
	}{
		{"empty value", "00", "", "0000"},
		{"single char", "01", "A", "0101A"},
		{"single digit length", "02", "hi", "0202hi"},
		{"double digit length", "03", "0123456789", "03100123456789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeTLV(tt.tag, tt.value)
			if got != tt.want {
				t.Errorf("encodeTLV(%q, %q): got %q, want %q", tt.tag, tt.value, got, tt.want)
			}
		})
	}
}

func Test_parseTLV(t *testing.T) {
	tests := []struct {
		name string
		data string
		want []tlvEntry
	}{
		{
			"single entry",
			"0103abc",
			[]tlvEntry{{Tag: "01", Value: "abc"}},
		},
		{
			"multiple entries",
			"0103abc" + "0202hi" + "0301X",
			[]tlvEntry{
				{Tag: "01", Value: "abc"},
				{Tag: "02", Value: "hi"},
				{Tag: "03", Value: "X"},
			},
		},
		{
			"empty value",
			"0100",
			[]tlvEntry{{Tag: "01", Value: ""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func Test_parseTLV_invalid(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{"truncated header", "01"},
		{"non-numeric length", "01XXa"},
		{"length exceeds data", "0105ab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseTLV(tt.data)
			if !errors.Is(err, ErrInvalidQR) {
				t.Errorf("got %v, want ErrInvalidQR", err)
			}
		})
	}
}
