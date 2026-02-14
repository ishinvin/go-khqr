package khqr

import "testing"

func TestCRC16Hex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data string
		want string
	}{
		{"empty", "", "FFFF"},
		{"single_byte", "A", "B915"},
		{"123456789", "123456789", "29B1"},
		{"hello", "hello", "D26E"},
		{"space", " ", "C592"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := crc16Hex(tt.data)
			if got != tt.want {
				t.Errorf("crc16Hex(%q) = %s, want %s", tt.data, got, tt.want)
			}
		})
	}
}

func TestCRC16HexUppercaseAndLength(t *testing.T) {
	t.Parallel()

	got := crc16Hex("hello")
	if len(got) != 4 {
		t.Fatalf("crc16Hex should return 4 chars, got %d", len(got))
	}
	for _, c := range got {
		if c >= 'a' && c <= 'f' {
			t.Errorf("crc16Hex should return uppercase hex, got %s", got)
			break
		}
	}
}
