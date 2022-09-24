package strings

import "fmt"

// eg: []byte{1, 0xA2, 0x33, 10} -> 01A2330A
func BytesAsLiteralString(b []byte) string {
	s := ""

	for i := 0; i < len(b); i++ {
		s += fmt.Sprintf("%02X", b[i])
	}

	return s
}
