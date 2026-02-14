package khqr

import "fmt"

// crcTable is a pre-computed CRC-16-CCITT lookup table (polynomial 0x1021).
var crcTable [256]uint16

func init() {
	for i := range 256 {
		crc := uint16(i) << 8 //nolint:gosec // i ranges [0,255], safe for uint16
		for range 8 {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc <<= 1
			}
		}
		crcTable[i] = crc
	}
}

// crc16 calculates a CRC-16-CCITT checksum with initial value 0xFFFF.
func crc16(data []byte) uint16 {
	crc := uint16(0xFFFF) //nolint:mnd // CRC-16-CCITT initial value
	for _, b := range data {
		crc = (crc << 8) ^ crcTable[byte(crc>>8)^b] //nolint:mnd // bits per byte
	}
	return crc
}

// crc16Hex computes the CRC-16-CCITT of the given string and returns it
// as a 4-character uppercase hexadecimal string.
func crc16Hex(data string) string {
	return fmt.Sprintf("%04X", crc16([]byte(data)))
}
